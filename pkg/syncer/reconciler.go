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
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"reflect"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	anyCluster = "*"
)

// SyncReconciler is the main structure that holds the state and dependencies for the multi-cluster syncer reconciler.
type SyncReconciler struct {
	storage storage.ResourceStorage

	client     client.Client
	controller controller.Controller
	mgr        MultiClusterSyncManager
}

// NewSyncReconciler creates a new instance of the SyncReconciler structure with the given storage.
func NewSyncReconciler(storage storage.ResourceStorage) *SyncReconciler {
	return &SyncReconciler{storage: storage}
}

// SetupWithManager sets up the SyncReconciler with the given manager and registers it as a controller.
func (r *SyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	controller, err := ctrl.NewControllerManagedBy(mgr).
		For(&clusterv1beta1.Cluster{}).
		Watches(&source.Kind{Type: &searchv1beta1.SyncRegistry{}}, &handler.Funcs{
			CreateFunc: r.CreateEvent,
			UpdateFunc: r.UpdateEvent,
			DeleteFunc: r.DeleteEvent,
		}).
		// TODO: watch syncResources & transformRule
		// Watches(&searchv1beta1.SyncResources{}).
		// Watches(&searchv1beta1.TransformRule{}).
		Build(r)
	if err != nil {
		return err
	}
	r.client = mgr.GetClient()
	r.controller = controller
	// TODO:
	r.mgr = NewMultiClusterSyncManager(context.Background(), r.controller, r.storage)
	return nil
}

// CreateEvent handles the creation event for a resource and enqueues it for reconciliation.
func (r *SyncReconciler) CreateEvent(ce event.CreateEvent, queue workqueue.RateLimitingInterface) {
	registry := ce.Object.(*searchv1beta1.SyncRegistry)
	for _, clusterName := range r.getMatchedClusters(registry) {
		queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{Name: clusterName}})
	}
}

// UpdateEvent handles the update event for a resource and enqueues it for reconciliation.
func (r *SyncReconciler) UpdateEvent(ue event.UpdateEvent, queue workqueue.RateLimitingInterface) {
	oldRegistry := ue.ObjectOld.(*searchv1beta1.SyncRegistry)
	newRegistry := ue.ObjectNew.(*searchv1beta1.SyncRegistry)

	if newRegistry.ResourceVersion == oldRegistry.ResourceVersion || reflect.DeepEqual(newRegistry.Spec, oldRegistry.Spec) {
		return
	}

	clusters := make(map[string]struct{})
	for _, clusterName := range r.getMatchedClusters(newRegistry) {
		clusters[clusterName] = struct{}{}
	}
	for _, clusterName := range r.getMatchedClusters(oldRegistry) {
		clusters[clusterName] = struct{}{}
	}
	for clusterName := range clusters {
		queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{Name: clusterName}})
	}
}

// DeleteEvent handles the deletion event for a resource and enqueues it for reconciliation.
func (r *SyncReconciler) DeleteEvent(de event.DeleteEvent, queue workqueue.RateLimitingInterface) {
	registry := de.Object.(*searchv1beta1.SyncRegistry)
	for _, clusterName := range r.getMatchedClusters(registry) {
		queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{Name: clusterName}})
	}
}

// Reconcile is the main entry point for the syncer reconciler, which is called whenever there is a change in the watched resources.
func (r *SyncReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
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
		return reconcile.Result{}, r.stopCluster(ctx, cluster.Name)
	}

	// TODO: it's danger
	// if !cluster.Status.Healthy {
	// 	logger.Info("cluster is unhealthy", "cluster", cluster.Name)
	// 	return reconcile.Result{}, r.stopCluster(ctx, cluster.Name)
	// }

	return reconcile.Result{}, r.handleClusterAddOrUpdate(ctx, cluster.DeepCopy())
}

// stopCluster stops the reconciliation process for the given cluster.
func (r *SyncReconciler) stopCluster(ctx context.Context, clusterName string) error {
	logger := ctrl.LoggerFrom(ctx)
	logger.Info("start to stop syncing cluster", "cluster", clusterName)
	if err := r.storage.DeleteAllResources(ctx, clusterName); err != nil {
		return err
	}
	r.mgr.Stop(ctx, clusterName)
	logger.Info("syncing cluster has been stopped", "cluster", clusterName)
	return nil
}

// startCluster starts the reconciliation process for the given cluster.
func (r *SyncReconciler) startCluster(ctx context.Context, clusterName string) error {
	logger := ctrl.LoggerFrom(ctx)

	logger.Info("start to sync cluster", "cluster", clusterName)
	return r.mgr.Start(ctx, clusterName)
}

// handleClusterAddOrUpdate is responsible for handling the addition or update of a cluster resource.
func (r *SyncReconciler) handleClusterAddOrUpdate(ctx context.Context, cluster *clusterv1beta1.Cluster) error {
	logger := ctrl.LoggerFrom(ctx)

	resources, err := r.getResources(ctx, cluster)
	if err != nil {
		return errors.Wrapf(err, "error detecting sync resources of the cluster %s", cluster.Name)
	}

	if len(resources) == 0 {
		logger.Info("cluster has no resources to sync", "cluster", cluster.Name)
		return r.stopCluster(ctx, cluster.Name)
	}

	clusterConfig, err := buildClusterConfig(cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to build config for cluster %s", cluster.Name)
	}

	singleMgr, exist := r.mgr.GetForCluster(cluster.Name)
	if exist && !reflect.DeepEqual(singleMgr.ClusterConfig(), clusterConfig) {
		logger.Info("cluster's spec has been changed, rebuild the manager", "cluster", cluster.Name)
		if err := r.stopCluster(ctx, cluster.Name); err != nil {
			return err
		}
		exist = false
	}
	if !exist {
		singleMgr, err = r.mgr.Create(ctx, cluster.Name, clusterConfig)
		if err != nil {
			return errors.Wrapf(err, "failed to setup the sync manager for cluster %s", cluster.Name)
		}

		if err := r.startCluster(ctx, cluster.Name); err != nil {
			return errors.Wrapf(err, "failed to start sync manager for cluster %s", err)
		}
	}

	if err := singleMgr.UpdateSyncResources(ctx, resources); err != nil {
		return errors.Wrapf(err, "failed to update sync resources for cluster %s", cluster.Name)
	}
	return nil
}

// getResources retrieves the list of resource sync rules for the given cluster.
func (r *SyncReconciler) getResources(ctx context.Context, cluster *clusterv1beta1.Cluster) ([]*searchv1beta1.ResourceSyncRule, error) {
	registries, err := r.getRegistries(ctx, cluster)
	if err != nil {
		return nil, err
	}

	var allResources []*searchv1beta1.ResourceSyncRule
	for _, registry := range registries {
		if !registry.DeletionTimestamp.IsZero() {
			continue
		}

		resources, err := r.getNormalizedResources(ctx, &registry)
		if err != nil {
			return nil, err
		}

		for _, r := range resources {
			allResources = append(allResources, r)
		}
	}
	return allResources, nil
}

// getNormalizedResources retrieves the normalized resource sync rules from the given sync registry.
func (r *SyncReconciler) getNormalizedResources(ctx context.Context, registry *searchv1beta1.SyncRegistry) (map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule, error) {
	var resources []searchv1beta1.ResourceSyncRule

	refName := registry.Spec.SyncResourcesRefName
	if refName != "" {
		var sr searchv1beta1.SyncResources
		err := r.client.Get(ctx, types.NamespacedName{Name: refName}, &sr)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get SyncResources %s", refName)
		}
		resources = append(resources, sr.Spec.SyncResources...)
	}

	resources = append(resources, registry.Spec.SyncResources...)

	ret := make(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	for _, res := range resources {
		nr, err := r.getNormalizedResource(ctx, &res)
		if err != nil {
			return nil, err
		}

		gvr, err := parseGVR(nr)
		if err != nil {
			return nil, err
		}

		ret[gvr] = nr
	}
	return ret, nil
}

// getRegistries retrieves the list of sync registries for the given cluster.
func (r *SyncReconciler) getRegistries(ctx context.Context, cluster *clusterv1beta1.Cluster) ([]searchv1beta1.SyncRegistry, error) {
	var syncRegistriesList searchv1beta1.SyncRegistryList
	if err := r.client.List(ctx, &syncRegistriesList); err != nil {
		return nil, errors.Wrap(err, "SyncRegistry lister error")
	}

	var ret []searchv1beta1.SyncRegistry
	for _, sr := range syncRegistriesList.Items {
		match, err := isMatched(&sr, cluster)
		if err != nil {
			return nil, err
		}
		if match {
			ret = append(ret, sr)
		}
	}
	return ret, nil
}

// getMatchedClusters returns the list of matched cluster names from the given sync registry.
func (r *SyncReconciler) getMatchedClusters(registry *searchv1beta1.SyncRegistry) []string {
	var ret []string

	clusters := &clusterv1beta1.ClusterList{}
	if err := r.client.List(context.Background(), clusters); err != nil {
		klog.Error("list clusters error: %v", err)
		return nil
	}

	for _, cluster := range clusters.Items {
		if match, _ := isMatched(registry, &cluster); match {
			ret = append(ret, cluster.Name)
		}
	}
	return ret
}

// isMatched checks if the given cluster matches the criteria specified in the sync registry.
func isMatched(registry *searchv1beta1.SyncRegistry, cluster *clusterv1beta1.Cluster) (bool, error) {
	for _, name := range registry.Spec.Clusters {
		if name == anyCluster || name == cluster.Name {
			return true, nil
		}
	}

	if registry.Spec.ClusterLabelSelector != nil {
		selector, err := metav1.LabelSelectorAsSelector(registry.Spec.ClusterLabelSelector)
		if err != nil {
			return false, err
		}
		if selector.Matches(labels.Set(cluster.GetLabels())) {
			return true, nil
		}
	}
	return false, nil
}

// getNormalizedResource retrieves the normalized resource sync rule for the given resource sync rule.
func (r *SyncReconciler) getNormalizedResource(ctx context.Context, rsr *searchv1beta1.ResourceSyncRule) (*searchv1beta1.ResourceSyncRule, error) {
	if rsr.TransformRefName == "" {
		return rsr.DeepCopy(), nil
	}
	if rsr.Transform != nil {
		return nil, fmt.Errorf("specify both Transform and TransformRefName in ResourceSyncRule is not allowed")
	}

	var rule searchv1beta1.TransformRule
	err := r.client.Get(ctx, types.NamespacedName{Name: rsr.TransformRefName}, &rule)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("TransformRule referenced by name %q doesn't exist", rsr.TransformRefName)
		}
		return nil, errors.Wrapf(err, "failed to list transformRule %s from lister", rsr.TransformRefName)
	}
	normalized := rsr.DeepCopy()
	normalized.Transform = &rule.Spec
	return normalized, nil
}

// buildClusterConfig creates a rest.Config object for the given cluster.
func buildClusterConfig(cluster *clusterv1beta1.Cluster) (*rest.Config, error) {
	access := cluster.Spec.Access
	if len(access.Endpoint) == 0 {
		return nil, fmt.Errorf("cluster %s's endpoint is empty", cluster.Name)
	}
	config := rest.Config{
		Host: access.Endpoint,
	}
	if access.Insecure != nil {
		config.TLSClientConfig.Insecure = *access.Insecure
	}
	if len(access.CABundle) > 0 {
		config.TLSClientConfig.CAData = access.CABundle
	}
	if access.Credential != nil {
		switch access.Credential.Type {
		case clusterv1beta1.CredentialTypeServiceAccountToken:
			config.BearerToken = access.Credential.ServiceAccountToken
		case clusterv1beta1.CredentialTypeX509Certificate:
			if access.Credential.X509 == nil {
				return nil, fmt.Errorf("cert and key is required for x509 credential type")
			}
			config.TLSClientConfig.CertData = access.Credential.X509.Certificate
			config.TLSClientConfig.KeyData = access.Credential.X509.PrivateKey
		case clusterv1beta1.CredentialTypeOIDC:
			if access.Credential.ExecConfig == nil {
				return nil, fmt.Errorf("ExecConfig is required for Exec credential type")
			}
			var env []clientcmdapi.ExecEnvVar
			for _, envValue := range access.Credential.ExecConfig.Env {
				tempEnv := clientcmdapi.ExecEnvVar{
					Name:  envValue.Name,
					Value: envValue.Value,
				}
				env = append(env, tempEnv)
			}
			config.ExecProvider = &clientcmdapi.ExecConfig{
				Command:            access.Credential.ExecConfig.Command,
				Args:               access.Credential.ExecConfig.Args,
				Env:                env,
				APIVersion:         access.Credential.ExecConfig.APIVersion,
				InstallHint:        access.Credential.ExecConfig.InstallHint,
				ProvideClusterInfo: access.Credential.ExecConfig.ProvideClusterInfo,
				InteractiveMode:    clientcmdapi.ExecInteractiveMode(access.Credential.ExecConfig.InteractiveMode),
			}
		default:
			return nil, fmt.Errorf("unknown credential type %v", access.Credential.Type)
		}
	}
	return &config, nil
}
