package search

import (
	"context"
	"fmt"
	"reflect"
	"time"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karbour/pkg/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/generated/informers/externalversions"
	clusterlisterv1beta1 "github.com/KusionStack/karbour/pkg/generated/listers/cluster/v1beta1"
	searchlisterv1beta1 "github.com/KusionStack/karbour/pkg/generated/listers/search/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type ClusterSyncController struct {
	clusterLister clusterlisterv1beta1.ClusterLister
	clusterSynced clientgocache.InformerSynced

	syncRegistryLister searchlisterv1beta1.SyncRegistryLister
	syncRegistrySynced clientgocache.InformerSynced

	syncResourcesLister searchlisterv1beta1.SyncResourcesLister
	syncResourcesSynced clientgocache.InformerSynced

	transformRuleLister searchlisterv1beta1.TransformRuleLister
	transformRuleSynced clientgocache.InformerSynced

	multiClusterSyncManager MultiClusterSyncManager
	queue                   workqueue.RateLimitingInterface
}

func NewSyncController(ctx context.Context, informerFactory externalversions.SharedInformerFactory, labelSelector labels.Selector) *ClusterSyncController {
	logger := klog.FromContext(ctx)

	clusterInformer := informerFactory.Cluster().V1beta1().Clusters()
	syncRegistryInformer := informerFactory.Search().V1beta1().SyncRegistries()
	syncResourcesInformer := informerFactory.Search().V1beta1().SyncResourceses()
	transformRuleInformer := informerFactory.Search().V1beta1().TransformRules()

	controller := &ClusterSyncController{
		clusterLister:       clusterInformer.Lister(),
		clusterSynced:       clusterInformer.Informer().HasSynced,
		syncRegistryLister:  syncRegistryInformer.Lister(),
		syncRegistrySynced:  syncRegistryInformer.Informer().HasSynced,
		syncResourcesLister: syncResourcesInformer.Lister(),
		syncResourcesSynced: syncResourcesInformer.Informer().HasSynced,
		transformRuleLister: transformRuleInformer.Lister(),
		transformRuleSynced: transformRuleInformer.Informer().HasSynced,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "sync-controller"),
	}
	logger.Info("Setting up event handlers")

	clusterInformer.Informer().AddEventHandler(clientgocache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addCluster,
		UpdateFunc: controller.updateCluster,
		DeleteFunc: controller.deleteCluster,
	})

	syncRegistryInformer.Informer().AddEventHandler(clientgocache.ResourceEventHandlerFuncs{
		AddFunc:    controller.addSyncRegistry,
		UpdateFunc: controller.updateSyncRegistry,
		DeleteFunc: controller.deleteSyncRegistry,
	})

	// TODO: syncResources event
	syncResourcesInformer.Informer().AddEventHandler(clientgocache.ResourceEventHandlerFuncs{})

	// TODO: transformRule event
	transformRuleInformer.Informer().AddEventHandler(clientgocache.ResourceEventHandlerFuncs{})

	return controller
}

func (c *ClusterSyncController) addCluster(obj interface{}) {
	cluster := obj.(*clusterv1beta1.Cluster)
	c.queue.Add(cluster.Name)
}

func (c *ClusterSyncController) updateCluster(oldObj interface{}, newObj interface{}) {
	oldCluster := oldObj.(*clusterv1beta1.Cluster)
	newCluster := newObj.(*clusterv1beta1.Cluster)
	if newCluster.ResourceVersion != oldCluster.ResourceVersion {
		c.queue.Add(newCluster.Name)
	}
}

func (c *ClusterSyncController) deleteCluster(obj interface{}) {
	cluster := obj.(*clusterv1beta1.Cluster)
	c.queue.Add(cluster.Name)
}

func (c *ClusterSyncController) addSyncRegistry(obj interface{}) {
	registry := obj.(*searchv1beta1.SyncRegistry)
	for _, clusterName := range c.getMatchedClusterNames(registry) {
		c.queue.Add(clusterName)
	}
}

func (c *ClusterSyncController) updateSyncRegistry(oldObj interface{}, newObj interface{}) {
	oldRegistry := oldObj.(*searchv1beta1.SyncRegistry)
	newRegistry := newObj.(*searchv1beta1.SyncRegistry)
	if newRegistry.ResourceVersion == oldRegistry.ResourceVersion || reflect.DeepEqual(newRegistry.Spec, oldRegistry.Spec) {
		return
	}

	clusters := make(map[string]struct{})
	for _, clusterName := range c.getMatchedClusterNames(newRegistry) {
		clusters[clusterName] = struct{}{}
	}
	for _, clusterName := range c.getMatchedClusterNames(oldRegistry) {
		clusters[clusterName] = struct{}{}
	}
	for clusterName := range clusters {
		c.queue.Add(clusterName)
	}
}

func (c *ClusterSyncController) deleteSyncRegistry(obj interface{}) {
	registry, ok := obj.(*searchv1beta1.SyncRegistry)
	if !ok {
		if tombstone, ok := obj.(clientgocache.DeletedFinalStateUnknown); ok {
			registry = tombstone.Obj.(*searchv1beta1.SyncRegistry)
		} else {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
	}

	for _, clusterName := range c.getMatchedClusterNames(registry) {
		c.queue.Add(clusterName)
	}
}

// TODO: record event
// TODO: update status
func (c *ClusterSyncController) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	logger := klog.FromContext(ctx)

	c.multiClusterSyncManager = NewMultiClusterSyncManager(ctx)

	// Start the informer factories to begin populating the informer caches
	logger.Info("Starting sync controller")

	// Wait for the caches to be synced before starting workers
	logger.Info("Waiting for informer caches to sync")

	if ok := clientgocache.WaitForCacheSync(ctx.Done(), c.clusterSynced, c.syncRegistrySynced, c.syncResourcesSynced, c.transformRuleSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	logger.Info("Starting workers", "count", workers)
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	logger.Info("Started workers")
	<-ctx.Done()
	logger.Info("Shutting down workers")

	return nil
}

func (c *ClusterSyncController) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *ClusterSyncController) processNextWorkItem(ctx context.Context) bool {
	obj, shutdown := c.queue.Get()
	logger := klog.FromContext(ctx)

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.queue.Done(obj)

		key, ok := obj.(string)
		if !ok {
			c.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.reconcile(ctx, key); err != nil {
			c.queue.AddRateLimited(key)
			return fmt.Errorf("error reconiling cluster %q: %s, requeuing", key, err.Error())
		}

		c.queue.Forget(obj)
		logger.Info("Successfully synced", "resource", key)
		return nil
	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true
}

func (c *ClusterSyncController) reconcile(ctx context.Context, clusterName string) error {
	logger := klog.FromContext(ctx)

	// TODO: handle local cluster case?

	cluster, err := c.clusterLister.Get(clusterName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("cluster doesn't exist", "cluster", clusterName)
			c.doRemoveCluster(ctx, clusterName)
			return nil
		}
		return err
	}

	if !cluster.DeletionTimestamp.IsZero() {
		logger.Info("cluster is being deleted", "cluster", clusterName)
		c.doRemoveCluster(ctx, clusterName)
		return nil
	}

	if !cluster.Status.Healthy {
		logger.Info("cluster is unhealthy", "cluster", clusterName)
		c.doRemoveCluster(ctx, clusterName)
		return nil
	}

	return c.doAddOrUpdateCluster(ctx, cluster.DeepCopy())
}

func (c *ClusterSyncController) doRemoveCluster(ctx context.Context, clusterName string) {
	logger := klog.FromContext(ctx)
	logger.Info("stopping syncing cluster", "cluster", clusterName)
	c.multiClusterSyncManager.Stop(clusterName)
	logger.Info("syncing cluster was stopped", "cluster", clusterName)
}

func (c *ClusterSyncController) doAddOrUpdateCluster(ctx context.Context, cluster *clusterv1beta1.Cluster) error {
	logger := klog.FromContext(ctx)

	resources, err := c.getClusterSyncResources(cluster)
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		logger.Info("cluster has no sync resources", "cluster", cluster.Name)
		c.doRemoveCluster(ctx, cluster.Name)
		return nil
	}

	config, err := buildClusterConfig(cluster)
	if err != nil {
		return err
	}
	mgr, exist := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
	if exist && !reflect.DeepEqual(mgr.ClusterConfig(), config) {
		logger.Info("cluster's spec has been changed, rebuild the sync manager", "cluster", cluster.Name)
		c.doRemoveCluster(ctx, cluster.Name)
		exist = false
	}
	if !exist {
		mgr, err = c.multiClusterSyncManager.ForCluster(cluster.Name, config)
		if err != nil {
			return fmt.Errorf("failed to build syncer for cluster %s: %v", cluster.Name, err)
		}
		if err := c.startCluster(ctx, cluster.Name); err != nil {
			return fmt.Errorf("failed to start syncer for cluster %s: %v", cluster.Name, err)
		}
	}
	mgr.SetSyncResources(resources)
	return nil
}

func (c *ClusterSyncController) startCluster(ctx context.Context, clusterName string) error {
	logger := klog.FromContext(ctx)
	logger.Info("start to sync cluster", "cluster", clusterName)
	return c.multiClusterSyncManager.Start(clusterName)
}

func buildClusterConfig(cluster *clusterv1beta1.Cluster) (*rest.Config, error) {
	access := cluster.Spec.Access
	if len(access.Endpoint) == 0 {
		return nil, fmt.Errorf("cluster %s's endpoint is empty", cluster.Name)
	}
	config := rest.Config{
		Host: access.Endpoint,
	}
	config.TLSClientConfig.Insecure = *access.Insecure
	if len(access.CABundle) > 0 {
		config.TLSClientConfig.CAData = access.CABundle
	}
	if access.Credential != nil {
		switch access.Credential.Type {
		case clusterv1beta1.CredentialTypeServiceAccountToken:
			// TODO: CredentialTypeServiceAccountToken
		case clusterv1beta1.CredentialTypeUnifiedIdentity:
			// TODO: CredentialTypeUnifiedIdentity
		case clusterv1beta1.CredentialTypeX509Certificate:
			if access.Credential.X509 == nil {
				return nil, fmt.Errorf("cert and key is required for x509 credential type")
			}
			config.TLSClientConfig.CertData = access.Credential.X509.Certificate
			config.TLSClientConfig.KeyData = access.Credential.X509.PrivateKey
		default:
			return nil, fmt.Errorf("unknown credential type %v", access.Credential.Type)
		}
	}
	return &config, nil
}

func (c *ClusterSyncController) getSyncRegistries(cluster *clusterv1beta1.Cluster) ([]*searchv1beta1.SyncRegistry, error) {
	syncRegistries, err := c.syncRegistryLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("SyncRegistry lister error: %v", err)
	}
	var ret []*searchv1beta1.SyncRegistry
	for _, sr := range syncRegistries {
		match, err := matchCluster(sr, cluster)
		if err != nil {
			return nil, err
		}
		if match {
			ret = append(ret, sr)
		}
	}
	return ret, nil
}

func (c *ClusterSyncController) getClusterSyncResources(cluster *clusterv1beta1.Cluster) (map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule, error) {
	registries, err := c.getSyncRegistries(cluster)
	if err != nil {
		return nil, err
	}

	ret := make(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	for _, registry := range registries {
		if !registry.DeletionTimestamp.IsZero() {
			continue
		}

		resources, err := c.getNormalizedSyncResources(registry)
		if err != nil {
			return nil, err
		}
		for gvr, rsr := range resources {
			if _, ok := ret[gvr]; ok {
				return nil, fmt.Errorf("found duplicate ResourceSyncRule defination for resource %q of cluster %q", gvr, cluster.Name)
			}
			ret[gvr] = rsr
		}
	}
	return ret, nil
}

func matchCluster(registry *searchv1beta1.SyncRegistry, cluster *clusterv1beta1.Cluster) (bool, error) {
	// TODO: cluster name & label selector 是AND的关系还是OR的关系
	for _, name := range registry.Spec.Clusters {
		if cluster.Name == name {
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

func (c *ClusterSyncController) getMatchedClusterNames(registry *searchv1beta1.SyncRegistry) []string {
	var ret []string
	clusters, _ := c.clusterLister.List(labels.Everything())
	for _, cluster := range clusters {
		if match, _ := matchCluster(registry, cluster); match {
			ret = append(ret, cluster.Name)
		}
	}
	return ret
}

func (c *ClusterSyncController) getNormalizedSyncResources(registry *searchv1beta1.SyncRegistry) (map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule, error) {
	var resources []searchv1beta1.ResourceSyncRule
	if registry.Spec.SyncResourcesRefName != "" {
		sr, err := c.syncResourcesLister.Get(registry.Spec.SyncResourcesRefName)
		if err != nil {
			return nil, fmt.Errorf("failed to get SyncResources %q: %v", registry.Spec.SyncResourcesRefName, err)
		}
		resources = append(resources, sr.Spec.SyncResources...)
	}
	resources = append(resources, registry.Spec.SyncResources...)

	ret := make(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	for _, r := range resources {
		nr, err := c.normalizeResourceSyncRule(&r)
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

func (c *ClusterSyncController) normalizeResourceSyncRule(rsr *searchv1beta1.ResourceSyncRule) (*searchv1beta1.ResourceSyncRule, error) {
	if rsr.TransformRefName == "" {
		return rsr, nil
	}
	if rsr.Transform != nil {
		return nil, fmt.Errorf("specify both Transform and TransformRefName in ResourceSyncRule is not allowed")
	}
	t, err := c.transformRuleLister.Get(rsr.TransformRefName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("TransformRule referenced by name %q doesn't exist", rsr.TransformRefName)
		}
		return nil, fmt.Errorf("failed to list transformRule %q from lister: %v", rsr.TransformRefName, err)
	}
	normalized := rsr.DeepCopy()
	normalized.Transform = &t.Spec
	return normalized, nil
}
