// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syncer

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sync"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
)

type filterForSdkFunc func(context.Context, schema.GroupVersionResource, *searchv1beta1.ResourceSyncRule) bool

type AgentReconciler struct {
	SyncReconciler
	gvrToGVKCache   sync.Map
	discoveryClient discovery.DiscoveryInterface
	clusterName     string
}

// NewAgentReconciler creates a new instance of the AgentReconciler structure with the given storage.
func NewAgentReconciler(storage storage.ResourceStorage, clusterName string) *AgentReconciler {
	return &AgentReconciler{
		SyncReconciler: SyncReconciler{
			storage: storage,
		},

		clusterName: clusterName,
	}
}

// SetupWithManager sets up the AgentReconciler with the given manager and registers it as a controller. Different from the SyncReconcile, it only focus on the special cluster.
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// only focus on the special cluster
	clusterFilter := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return e.Object.GetName() == r.clusterName
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.GetName() == r.clusterName
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return e.Object.GetName() == r.clusterName
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return e.Object.GetName() == r.clusterName
		},
	}

	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1beta1.Cluster{}, builder.WithPredicates(clusterFilter)).
		Watches(&source.Kind{Type: &searchv1beta1.SyncRegistry{}}, &handler.Funcs{
			CreateFunc: r.CreateEvent,
			UpdateFunc: r.UpdateEvent,
			DeleteFunc: r.DeleteEvent,
		}).
		Build(r)
	if err != nil {
		return err
	}
	r.client = mgr.GetClient()
	r.controller = controller
	r.discoveryClient = discovery.NewDiscoveryClientForConfigOrDie(mgr.GetConfig())
	// TODO:
	r.mgr = NewMultiClusterSyncManager(context.Background(), r.controller, r.storage)
	return nil
}

// Reconcile is the main entry point for the syncer reconciler, which is called whenever there is a change in the watched resources.
func (r *AgentReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	var cluster clusterv1beta1.Cluster
	if err := r.client.Get(ctx, req.NamespacedName, &cluster); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("cluster doesn't exist", "cluster", req.Name)
			return reconcile.Result{}, r.stopCluster(ctx, req.Name)
		}
		return reconcile.Result{}, err
	}

	if !cluster.DeletionTimestamp.IsZero() {
		logger.Info("cluster is being deleted", "cluster", cluster.Name)

		err := r.stopCluster(ctx, cluster.Name)
		if err != nil {
			return reconcile.Result{}, err
		}

		if len(cluster.Finalizers) > 0 {
			cluster.Finalizers = nil
			err := r.client.Update(ctx, &cluster)
			if err != nil && !apierrors.IsNotFound(err) {
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, r.handleClusterAddOrUpdate(ctx, cluster.DeepCopy(), buildClusterConfigInAgent, r.filterForSdk)
}

func (r *AgentReconciler) filterForSdk(ctx context.Context, gvr schema.GroupVersionResource, rule *searchv1beta1.ResourceSyncRule) bool {
	gvk, err := r.gvrToGVK(ctx, gvr)
	if err != nil {
		return false
	}
	if _, exist := utils.GetSyncGVK(gvk); exist {
		// If gvk map consist special gvk, it means that gvk is already reconciled by dynamic reconciler
		utils.SetSyncGVK(gvk, *rule)
		return false
	}
	return true
}

func (r *AgentReconciler) gvrToGVK(ctx context.Context, gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	if val, ok := r.gvrToGVKCache.Load(gvr); ok {
		return val.(schema.GroupVersionKind), nil
	}

	zero := schema.GroupVersionKind{}

	groupResources, err := restmapper.GetAPIGroupResources(r.discoveryClient)
	if err != nil {
		return zero, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)
	gvk, err := mapper.KindFor(gvr)
	if err != nil {
		return zero, err
	}

	ctrl.LoggerFrom(ctx).Info("GVR to GVV", "gvr", gvr, "gvk", gvk)
	r.gvrToGVKCache.Store(gvr, gvk)
	return gvk, nil
}

func buildClusterConfigInAgent(cluster *clusterv1beta1.Cluster) (*rest.Config, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{
		WarnIfAllMissing: false,
		Precedence:       []string{clientcmd.RecommendedHomeFile},
		MigrationRules: map[string]string{
			clientcmd.RecommendedHomeFile: filepath.Join(os.Getenv("HOME"), clientcmd.RecommendedHomeDir, ".kubeconfig"),
		},
	}
	if _, ok := os.LookupEnv("HOME"); !ok {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("could not get current user: %v", err)
		}
		loadingRules.Precedence = append(loadingRules.Precedence, path.Join(u.HomeDir, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))
	}
	cfg, err := loadConfigWithContext("", loadingRules, "")
	if err != nil {
		return nil, err
	}
	if cfg.QPS == 0.0 {
		cfg.QPS = 20.0
		cfg.Burst = 30.0
	}
	return cfg, nil
}

func loadConfigWithContext(apiServerURL string, loader clientcmd.ClientConfigLoader, context string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loader,
		&clientcmd.ConfigOverrides{
			ClusterInfo: clientcmdapi.Cluster{
				Server: apiServerURL,
			},
			CurrentContext: context,
		}).ClientConfig()
}
