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
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"reflect"
	templateUtil "text/template"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilErr "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/KusionStack/karpor/config"
	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	clusterv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	"github.com/KusionStack/karpor/pkg/syncer/template"
	"github.com/KusionStack/karpor/pkg/syncer/utils"
	"github.com/KusionStack/karpor/pkg/util/certgenerator"

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

	highAvailability bool

	client     client.Client
	controller controller.Controller
	mgr        MultiClusterSyncManager

	storageAddresses []string

	externalEndpoint string

	caCert *x509.Certificate
	caKey  crypto.Signer
}

// NewSyncReconciler creates a new instance of the SyncReconciler structure with the given storage.
func NewSyncReconciler(storage storage.ResourceStorage, highAvailability bool, storageAddresses []string,
	externalEndpoint string, caCert *x509.Certificate, caKey crypto.Signer) *SyncReconciler {
	return &SyncReconciler{
		storage:          storage,
		highAvailability: highAvailability,
		storageAddresses: storageAddresses,
		externalEndpoint: externalEndpoint,
		caCert:           caCert,
		caKey:            caKey,
	}
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

	if r.highAvailability {
		if cluster.Spec.Mode == clusterv1beta1.PushClusterMode {
			return reconcile.Result{}, r.handleClusterAddOrUpdateForPush(ctx, cluster.DeepCopy())
		}
		// default pull mode
		return reconcile.Result{}, r.handleClusterAddOrUpdateForPull(ctx, cluster.DeepCopy())
	}

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

// handleClusterAddOrUpdateForPush dispatches the relevant crds resources to user cluster for a push mode cluster resource in high-availability scene.
func (r *SyncReconciler) handleClusterAddOrUpdateForPush(ctx context.Context, cluster *clusterv1beta1.Cluster) error {
	logger := ctrl.LoggerFrom(ctx)
	logger.V(5).Info("handle cluster has been added/updated, push mode", "cluster", cluster.Name)

	// build user client
	clusterConfig, err := buildClusterConfig(cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to build config for cluster %s", cluster.Name)
	}
	dynamicClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return errors.Wrapf(err, "failed to build dynamic client for cluster %s", cluster.Name)
	}

	// must apply crds before other resources.
	err = utils.ApplyCrds(ctx, dynamicClient)
	if err != nil {
		return errors.Wrapf(err, "failed to apply crds for cluster %s", cluster.Name)
	}

	err = r.dispatchResources(ctx, dynamicClient, cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to dispatch resources for cluster %s", cluster.Name)
	}

	err = r.generateAgentYaml(ctx, cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to generate agent yaml for push mode cluster %s", cluster.Name)
	}

	return err
}

// handleClusterAddOrUpdate generate the cert for agent in user cluster
func (r *SyncReconciler) handleClusterAddOrUpdateForPull(ctx context.Context, cluster *clusterv1beta1.Cluster) error {
	err := r.createClusterRoleAndBinding(ctx, cluster)
	if err != nil {
		return err
	}

	err = r.generateAgentYaml(ctx, cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to generate agent yaml for pull mmode cluster %s", cluster.Name)
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

// dispatchResources dispatch cluster, syncregistry, syncresource, transformrule and trimrule resource to user cluster.
func (r *SyncReconciler) dispatchResources(ctx context.Context, dynamicClient dynamic.Interface, cluster *clusterv1beta1.Cluster) error {
	// collect the resources needed to be dispatched
	unstructuredObjectMap := map[schema.GroupVersionResource][]unstructured.Unstructured{}
	err := r.getUnstructuredCluster(cluster, unstructuredObjectMap)
	if err != nil {
		return errors.Wrapf(err, "error get unstructured cluster cr for cluster %s", cluster.Name)
	}

	err = r.getUnstructuredRegistries(ctx, cluster, unstructuredObjectMap)
	if err != nil {
		return errors.Wrapf(err, "error get unstructured objects of the syncregistries for cluster %s", cluster.Name)
	}

	// dispatch resources
	var errs []error
	for gvr := range unstructuredObjectMap {
		for idx := range unstructuredObjectMap[gvr] {
			unstructuredObj := unstructuredObjectMap[gvr][idx]
			err := utils.CreateOrUpdateUnstructured(ctx, dynamicClient, gvr, "", &unstructuredObj)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilErr.NewAggregate(errs)
}

// createClusterRoleAndBinding create clusterrole and clusterrolebinding for agent in pull mode
func (r *SyncReconciler) createClusterRoleAndBinding(ctx context.Context, cluster *clusterv1beta1.Cluster) error {
	clusterRoleName := fmt.Sprintf("%s-clusterrole", cluster.Name)
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{clusterv1beta1.SchemeGroupVersion.Group},
				Resources:     []string{"clusters"},
				Verbs:         []string{"get"},
				ResourceNames: []string{cluster.Name},
			},
			{
				APIGroups: []string{clusterv1beta1.SchemeGroupVersion.Group},
				Resources: []string{"clusters"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{searchv1beta1.SchemeGroupVersion.Group},
				Resources: []string{"syncregistries"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{searchv1beta1.SchemeGroupVersion.Group},
				Resources: []string{"syncclusterresources"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{searchv1beta1.SchemeGroupVersion.Group},
				Resources: []string{"transformrules"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{searchv1beta1.SchemeGroupVersion.Group},
				Resources: []string{"trimrules"},
				Verbs:     []string{"get", "watch", "list"},
			},
		},
	}

	err := r.client.Create(ctx, clusterRole)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "failed to create clusterrole %s", clusterRole.Name)
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-clusterrolebinding", cluster.Name),
		},
		Subjects: []rbacv1.Subject{
			{
				APIGroup: rbacv1.SchemeGroupVersion.Group,
				Kind:     rbacv1.UserKind,
				Name:     cluster.Name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
	}

	err = r.client.Create(ctx, clusterRoleBinding)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return errors.Wrapf(err, "failed to create clusterrolebinding %s", clusterRoleBinding.Name)
	}

	return nil
}

// generateAgentYaml generate the agent yaml to be deployed in user cluster, and save yaml into secret.
func (r *SyncReconciler) generateAgentYaml(ctx context.Context, cluster *clusterv1beta1.Cluster) error {
	var certData, keyData []byte
	if cluster.Spec.Mode == clusterv1beta1.PullClusterMode {
		certConfig := certgenerator.Config{
			CAName:     "kubernetes",
			CommonName: cluster.Name,
			Year:       100,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		}
		cert, key, err := certgenerator.NewCaCertAndKeyFromRoot(certConfig, r.caCert, r.caKey)
		if err != nil {
			return err
		}
		certData = certgenerator.EncodeCertPEM(cert)
		keyData, err = keyutil.MarshalPrivateKeyToPEM(key)
		if err != nil {
			return fmt.Errorf("unable to marshal private key to PEM %s", err)
		}
	}

	agentYml, err := renderYamlFile(cluster, r.storageAddresses, r.externalEndpoint,
		base64.StdEncoding.EncodeToString(certData), base64.StdEncoding.EncodeToString(keyData))
	if err != nil {
		return errors.Wrap(err, "failed to render agent yaml")
	}

	err = applyAgentYmlSecret(ctx, r.client, cluster, []byte(agentYml))
	if err != nil {
		return err
	}

	return nil
}

// getUnstructuredCluster retrieves the cluster cr for the given cluster.
func (r *SyncReconciler) getUnstructuredCluster(cluster *clusterv1beta1.Cluster, unstructuredObjectMap map[schema.GroupVersionResource][]unstructured.Unstructured) error {
	// only dispatch cluster cr once
	unstructuredObj, err := utils.ConvertToUnstructured(cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to convert to unstructured object")
	}
	unstructuredObjectMap[clusterv1beta1.SchemeGroupVersion.WithResource("clusters")] = []unstructured.Unstructured{*unstructuredObj}

	return nil
}

// getRegistries retrieves the list of sync registries for the given cluster.
func (r *SyncReconciler) getUnstructuredRegistries(ctx context.Context, cluster *clusterv1beta1.Cluster, unstructuredObjectMap map[schema.GroupVersionResource][]unstructured.Unstructured) error {
	registries, err := r.getRegistries(ctx, cluster)
	if err != nil {
		return errors.Wrapf(err, "failed to get registries")
	}

	// init map value
	unstructuredRegistries := make([]unstructured.Unstructured, 0, len(registries))
	var unstructuredTransformRules []unstructured.Unstructured
	var unstructuredTrimRules []unstructured.Unstructured

	// avoid to collect duplicate cr
	transformRuleMap := make(map[string]struct{})
	trimRuleMap := make(map[string]struct{})

	// collect cr list
	for idx := range registries {
		registry := registries[idx]
		// set special cluster name when dispatching
		registry.Spec.Clusters = []string{cluster.Name}

		unstructuredObj, err := utils.ConvertToUnstructured(&registry)
		if err != nil {
			return errors.Wrapf(err, "failed to convert to unstructured object for registry %s", registry.Name)
		}
		// do not set status when update
		unstructuredObj.Object["status"] = nil
		unstructuredRegistries = append(unstructuredRegistries, *unstructuredObj)

		// obtain relevant cr
		for _, sr := range registry.Spec.SyncResources {
			if _, ok := transformRuleMap[sr.TransformRefName]; !ok && sr.TransformRefName != "" {
				rule, err := r.extractTransformRule(ctx, &sr)
				if err != nil {
					return err
				}
				unstructuredObj, err = utils.ConvertToUnstructured(rule)
				if err != nil {
					return errors.Wrapf(err, "failed to convert to unstructured object for transformrule %s", rule.Name)
				}

				unstructuredTransformRules = append(unstructuredTransformRules, *unstructuredObj)
				transformRuleMap[sr.TransformRefName] = struct{}{}
			}

			if _, ok := trimRuleMap[sr.TrimRefName]; !ok && sr.TrimRefName != "" {
				rule, err := r.extractTrimRule(ctx, &sr)
				if err != nil {
					return err
				}
				unstructuredObj, err = utils.ConvertToUnstructured(rule)
				if err != nil {
					return errors.Wrapf(err, "failed to convert to unstructured object for trimrule %s", rule.Name)
				}

				unstructuredTrimRules = append(unstructuredTrimRules, *unstructuredObj)
				transformRuleMap[sr.TrimRefName] = struct{}{}
			}
		}
	}

	// set map
	unstructuredObjectMap[searchv1beta1.SchemeGroupVersion.WithResource("syncregistries")] = unstructuredRegistries
	unstructuredObjectMap[searchv1beta1.SchemeGroupVersion.WithResource("transformrules")] = unstructuredTransformRules
	unstructuredObjectMap[searchv1beta1.SchemeGroupVersion.WithResource("trimrules")] = unstructuredTrimRules

	return nil
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
	normalized := rsr.DeepCopy()

	if rsr.TransformRefName != "" {
		rule, err := r.extractTransformRule(ctx, rsr)
		if err != nil {
			return nil, err
		}
		normalized.Transform = &rule.Spec
	}

	if rsr.TrimRefName != "" {
		rule, err := r.extractTrimRule(ctx, rsr)
		if err != nil {
			return nil, err
		}
		normalized.Trim = &rule.Spec
	}

	return normalized, nil
}

func (r *SyncReconciler) extractTransformRule(ctx context.Context, rsr *searchv1beta1.ResourceSyncRule) (*searchv1beta1.TransformRule, error) {
	var rule searchv1beta1.TransformRule
	if rsr.Transform != nil {
		return nil, fmt.Errorf("specify both Transform and TransformRefName in ResourceSyncRule is not allowed")
	}

	err := r.client.Get(ctx, types.NamespacedName{Name: rsr.TransformRefName}, &rule)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("TransformRule referenced by name %q doesn't exist", rsr.TransformRefName)
		}
		return nil, errors.Wrapf(err, "failed to list transformRule %s from lister", rsr.TransformRefName)
	}
	return &rule, nil
}

func (r *SyncReconciler) extractTrimRule(ctx context.Context, rsr *searchv1beta1.ResourceSyncRule) (*searchv1beta1.TrimRule, error) {
	var rule searchv1beta1.TrimRule
	if rsr.Trim != nil {
		return nil, fmt.Errorf("specify both Trim and TrimRefName in ResourceSyncRule is not allowed")
	}

	err := r.client.Get(ctx, types.NamespacedName{Name: rsr.TrimRefName}, &rule)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("TrimRule referenced by name %q doesn't exist", rsr.TrimRefName)
		}
		return nil, errors.Wrapf(err, "failed to list trimRule %s from lister", rsr.TrimRefName)
	}
	return &rule, nil
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

func renderYamlFile(cluster *clusterv1beta1.Cluster, storageAddresses []string, externalEndpoint, certData, keyData string) (string, error) {
	c := template.Config{
		Mode:             cluster.Spec.Mode,
		ClusterName:      cluster.Name,
		Level:            cluster.Spec.Level,
		StorageAddresses: storageAddresses,
		ClusterMode:      cluster.Spec.Mode,
	}

	if cluster.Spec.Mode == clusterv1beta1.PullClusterMode {
		c.ExternalEndpoint = externalEndpoint
		c.CaCert = certData
		c.CaKey = keyData
	}

	agentYml, err := templateUtil.New("").Parse(string(config.AgentTpl))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse agent yaml")
	}

	var renderedTemplate bytes.Buffer
	err = agentYml.Execute(&renderedTemplate, c)
	if err != nil {
		return "", errors.Wrap(err, "failed to render agent yaml")
	}

	return renderedTemplate.String(), nil
}

func applyAgentYmlSecret(ctx context.Context, cli client.Client, cluster *clusterv1beta1.Cluster, content []byte) error {
	secretName := fmt.Sprintf("%s-agent", cluster.Name)
	newAgentSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: "karpor",
		},
		Data: map[string][]byte{
			"config": content,
		},
	}

	oldAgentSecret := &v1.Secret{}
	err := cli.Get(ctx, client.ObjectKey{
		Name:      secretName,
		Namespace: "karpor",
	}, oldAgentSecret)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return errors.Wrap(err, "failed to create agent secret")
		}
		err = cli.Create(ctx, newAgentSecret)
		if err != nil {
			return errors.Wrap(err, "failed to update agent secret")
		}
	} else if !reflect.DeepEqual(oldAgentSecret.Data, newAgentSecret.Data) {
		err = cli.Update(ctx, newAgentSecret)
		if err != nil {
			return errors.Wrap(err, "failed to update agent secret")
		}
	}
	return nil
}
