package search

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	clusterv1beta1 "github.com/KusionStack/karbour/pkg/apis/cluster/v1beta1"
	searchv1beta1 "github.com/KusionStack/karbour/pkg/apis/search/v1beta1"
	"github.com/KusionStack/karbour/pkg/generated/clientset/versioned/fake"
	"github.com/KusionStack/karbour/pkg/generated/informers/externalversions"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	clientgocache "k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/ktesting"
)

var (
	alwaysReady        = func() bool { return true }
	noResyncPeriodFunc = func() time.Duration { return 0 }
)

type fixture struct {
	t *testing.T

	client *fake.Clientset
	// Objects to put in the store.
	clusterLister       []*clusterv1beta1.Cluster
	syncRegistryLister  []*searchv1beta1.SyncRegistry
	syncResourcesLister []*searchv1beta1.SyncResources
	transformRuleLister []*searchv1beta1.TransformRule
	// Objects from here preloaded into NewSimpleFake.
	objects []runtime.Object
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.objects = []runtime.Object{}
	return f
}

func (f *fixture) newSyncController(ctx context.Context) (*ClusterSyncController, externalversions.SharedInformerFactory) {
	f.client = fake.NewSimpleClientset(f.objects...)
	factory := externalversions.NewSharedInformerFactory(f.client, noResyncPeriodFunc())
	controller := NewSyncController(ctx, factory, nil)

	// controller.clusterSynced = alwaysReady
	// controller.syncRegistrySynced = alwaysReady
	// controller.syncResourcesSynced = alwaysReady
	// controller.transformRuleSynced = alwaysReady

	for _, cluster := range f.clusterLister {
		factory.Cluster().V1beta1().Clusters().Informer().GetIndexer().Add(cluster)
	}
	for _, syncRegistry := range f.syncRegistryLister {
		factory.Search().V1beta1().SyncRegistries().Informer().GetIndexer().Add(syncRegistry)
	}
	for _, syncResources := range f.syncResourcesLister {
		factory.Search().V1beta1().SyncResourceses().Informer().GetIndexer().Add(syncResources)
	}
	for _, transformRule := range f.transformRuleLister {
		factory.Search().V1beta1().TransformRules().Informer().GetIndexer().Add(transformRule)
	}
	return controller, factory
}

func (f *fixture) run(ctx context.Context, clusterName string) {
	f.runReconcile(ctx, clusterName, true)
}

func (f *fixture) runReconcile(ctx context.Context, clusterName string, startInformers bool) {
	c, factory := f.newSyncController(ctx)
	if startInformers {
		factory.Start(ctx.Done())
	}

	if c.multiClusterSyncManager == nil {
		c.multiClusterSyncManager = NewMultiClusterSyncManager(ctx)
	}

	err := c.reconcile(ctx, clusterName)
	if err != nil {
		f.t.Errorf("error syncing cluster: %v", err)
	}

	mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(clusterName)
	assert.True(f.t, ok)
	assert.NotNil(f.t, mgr)
	assert.True(f.t, mgr.Started())
}

func newCluster(name string, labels ...string) *clusterv1beta1.Cluster {
	lables_mp := make(map[string]string)
	for _, label := range labels {
		parts := strings.Split(label, "/")
		lables_mp[parts[0]] = parts[1]
	}

	return &clusterv1beta1.Cluster{
		TypeMeta: metav1.TypeMeta{APIVersion: clusterv1beta1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: lables_mp,
		},
		Status: clusterv1beta1.ClusterStatus{
			Healthy: true,
		},
	}
}

func newSyncRegistry(name string, clusterNames []string, syncResourcesRefName string) *searchv1beta1.SyncRegistry {

	return &searchv1beta1.SyncRegistry{
		TypeMeta: metav1.TypeMeta{APIVersion: searchv1beta1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: searchv1beta1.SyncRegistrySpec{
			Clusters:             clusterNames,
			SyncResourcesRefName: syncResourcesRefName,
		},
	}
}

func newSyncResource(name string, resourceSyncRules []searchv1beta1.ResourceSyncRule) *searchv1beta1.SyncResources {
	return &searchv1beta1.SyncResources{
		TypeMeta: metav1.TypeMeta{APIVersion: searchv1beta1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: searchv1beta1.SyncResourcesSpec{
			SyncResources: resourceSyncRules,
		},
	}
}

func newResourceSyncRule(APIVersion string, resource string, namespace string) searchv1beta1.ResourceSyncRule {
	return searchv1beta1.ResourceSyncRule{
		APIVersion: APIVersion,
		Resource:   resource,
		Namespace:  namespace,
	}
}

func newTransformRule(name string, tp string, valueTemplate string) *searchv1beta1.TransformRule {
	return &searchv1beta1.TransformRule{
		TypeMeta: metav1.TypeMeta{APIVersion: searchv1beta1.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: searchv1beta1.TransformRuleSpec{
			Type:          tp,
			ValueTemplate: valueTemplate,
		},
	}
}

type fakeMultiClusterSyncManager struct {
	syncers map[string]SingleClusterSyncManager
	sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	logger klog.Logger
}

func newFakeMultiClusterSyncManager(baseContext context.Context) MultiClusterSyncManager {
	innerCtx, innerCancel := context.WithCancel(baseContext)
	return &fakeMultiClusterSyncManager{
		syncers: make(map[string]SingleClusterSyncManager),
		ctx:     innerCtx,
		cancel:  innerCancel,
		logger:  klog.FromContext(baseContext),
	}
}

func (s *fakeMultiClusterSyncManager) ForCluster(clusterName string, config *rest.Config) (SingleClusterSyncManager, error) {
	mgr, ok := s.GetSingleClusterSyncManager(clusterName)
	if ok {
		// already exist, just return
		return mgr, nil
	}

	mgr, err := newFakeSingleClusterSyncManager(s.ctx, clusterName, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build sync manager for cluster %q: %v", clusterName, err)
	}
	s.Lock()
	defer s.Unlock()
	s.syncers[clusterName] = mgr
	return mgr, nil
}

func (s *fakeMultiClusterSyncManager) Start(clusterName string) error {
	mgr, ok := s.GetSingleClusterSyncManager(clusterName)
	if !ok {
		return fmt.Errorf("SingleClusterSyncManager for cluster %q does not exist", clusterName)
	}
	return mgr.Start()
}

func (s *fakeMultiClusterSyncManager) Stop(clusterName string) {
	mgr, found := s.GetSingleClusterSyncManager(clusterName)
	if !found {
		return
	}
	mgr.Stop()

	s.Lock()
	defer s.Unlock()
	delete(s.syncers, clusterName)
}

func (s *fakeMultiClusterSyncManager) GetSingleClusterSyncManager(clusterName string) (SingleClusterSyncManager, bool) {
	s.RLock()
	defer s.RUnlock()
	syncer, ok := s.syncers[clusterName]
	return syncer, ok
}

type fakeSingleClusterSyncManager struct {
	clusterName   string
	clusterConfig *rest.Config

	started, stopped bool
	syncResources    atomic.Value
	logger           klog.Logger
}

func newFakeSingleClusterSyncManager(baseContext context.Context, clusterName string, config *rest.Config) (SingleClusterSyncManager, error) {
	config = rest.CopyConfig(config)

	mgr := &fakeSingleClusterSyncManager{
		clusterName:   clusterName,
		clusterConfig: config,
		logger:        klog.FromContext(baseContext).WithValues("cluster", clusterName),
	}
	return mgr, nil
}

func (s *fakeSingleClusterSyncManager) Started() bool {
	return s.started
}

func (s *fakeSingleClusterSyncManager) Stopped() bool {
	return s.stopped
}

func (s *fakeSingleClusterSyncManager) Start() error {
	s.logger.Info("start sync manager")
	s.started = true
	s.process()
	return nil
}

func (s *fakeSingleClusterSyncManager) Stop() {
	s.logger.Info("stop sync manager")
	s.stopped = true
}

func (s *fakeSingleClusterSyncManager) SetSyncResources(newSyncResources map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule) {
	s.logger.Info("update sync resources")
	s.syncResources.Store(newSyncResources)

}

func (s *fakeSingleClusterSyncManager) process() {
	s.logger.Info("start sync manager process")
}

func (s *fakeSingleClusterSyncManager) IsSyncerExist(gvr schema.GroupVersionResource) bool {
	syncResources := s.syncResources.Load().(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	_, ok := syncResources[gvr]
	return ok
}

func (s *fakeSingleClusterSyncManager) GetSyncResources() map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule {
	syncResources := s.syncResources.Load().(map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule)
	return syncResources
}

func (s *fakeSingleClusterSyncManager) ClusterConfig() *rest.Config {
	return rest.CopyConfig(s.clusterConfig)
}

// run with fake sync manager
func (c *ClusterSyncController) run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	logger := klog.FromContext(ctx)

	c.multiClusterSyncManager = newFakeMultiClusterSyncManager(ctx)

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

func TestReconcile(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, cluster, syncRegistry)

	f.run(ctx, "test-cluster")
}

func TestCreateCluster(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, syncRegistry)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.ClusterV1beta1().Clusters().Create(ctx, cluster, metav1.CreateOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	mgr, _ := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
	assert.True(t, mgr.Started())
	assert.True(t, mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}))
}

func TestUpdateCluster(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, syncRegistry)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.ClusterV1beta1().Clusters().Create(ctx, cluster, metav1.CreateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	cluster, err = factory.Cluster().V1beta1().Clusters().Lister().Get(cluster.Name)
	assert.NoError(t, err)
	clusterCopy := cluster.DeepCopy()
	cluster.ResourceVersion = "1"
	clusterCopy.Spec.Access.Endpoint = "https://test-cluster2"
	_, err = f.client.ClusterV1beta1().Clusters().Update(ctx, clusterCopy, metav1.UpdateOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		if mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(clusterCopy.Name); ok {
			return mgr.ClusterConfig().Host == clusterCopy.Spec.Access.Endpoint
		}
		return false
	}, wait.ForeverTestTimeout, 100*time.Millisecond)
}

func TestDeleteCluster(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, syncRegistry)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.ClusterV1beta1().Clusters().Create(ctx, cluster, metav1.CreateOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	err = f.client.ClusterV1beta1().Clusters().Delete(ctx, cluster.Name, metav1.DeleteOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return !ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)
}

func TestCreateSyncRegistry(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.objects = append(f.objects, cluster)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.SearchV1beta1().SyncRegistries().Create(ctx, syncRegistry, metav1.CreateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	mgr, _ := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
	assert.True(t, mgr.Started())
	assert.True(t, mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}))

}

func TestUpdateSyncRegistry(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, cluster, syncRegistry)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	syncRegistry, err := f.client.SearchV1beta1().SyncRegistries().Get(ctx, syncRegistry.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	syncRegistryCopy := syncRegistry.DeepCopy()
	syncRegistryCopy.Spec.SyncResources = append(syncRegistryCopy.Spec.SyncResources, newResourceSyncRule("v1", "services", "default"))
	syncRegistryCopy.ResourceVersion = "1"
	_, err = f.client.SearchV1beta1().SyncRegistries().Update(ctx, syncRegistryCopy, metav1.UpdateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		if mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name); !ok {
			return false
		} else {
			return mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}) && mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"})
		}
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

}

func TestDeleteSyncRegistry(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRegistry.Spec.SyncResources = []searchv1beta1.ResourceSyncRule{syncRule}

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.objects = append(f.objects, cluster)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.SearchV1beta1().SyncRegistries().Create(ctx, syncRegistry, metav1.CreateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	err = f.client.SearchV1beta1().SyncRegistries().Delete(ctx, syncRegistry.Name, metav1.DeleteOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return !ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)
}

func TestCreateSyncResources(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "test-sync-resource")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncResource := newSyncResource("test-sync-resource", []searchv1beta1.ResourceSyncRule{syncRule})

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, cluster, syncRegistry)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.SearchV1beta1().SyncResourceses().Create(ctx, syncResource, metav1.CreateOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		if mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name); !ok {
			return false
		} else {
			return mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"})
		}
	}, wait.ForeverTestTimeout, 100*time.Millisecond)
}

func TestUpdateSyncResources(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "test-sync-resource")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncResource := newSyncResource("test-sync-resource", []searchv1beta1.ResourceSyncRule{syncRule})

	_, ctx := ktesting.NewTestContext(t)
	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.syncResourcesLister = append(f.syncResourcesLister, syncResource)
	f.objects = append(f.objects, cluster, syncRegistry, syncResource)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	syncResource, err := f.client.SearchV1beta1().SyncResourceses().Get(ctx, syncResource.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	syncResourceCopy := syncResource.DeepCopy()
	syncResourceCopy.Spec.SyncResources[0].APIVersion = "v2"
	syncResourceCopy.Spec.SyncResources[0].Namespace = "kube-system"
	syncResourceCopy.Spec.SyncResources[0].Resource = "configmaps"
	syncResourceCopy.ResourceVersion = "1"
	_, err = f.client.SearchV1beta1().SyncResourceses().Update(ctx, syncResourceCopy, metav1.UpdateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		if mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name); !ok {
			return false
		} else {
			return mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v2", Resource: "configmaps"})
		}
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

}

func TestDeleteSyncResources(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "test-sync-resource")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncResource := newSyncResource("test-sync-resource", []searchv1beta1.ResourceSyncRule{syncRule})

	_, ctx := ktesting.NewTestContext(t)
	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.objects = append(f.objects, cluster, syncRegistry)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.Run(ctx, 1)

	_, err := f.client.SearchV1beta1().SyncResourceses().Create(ctx, syncResource, metav1.CreateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		if mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name); !ok {
			return false
		} else {
			return mgr.IsSyncerExist(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"})
		}
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	err = f.client.SearchV1beta1().SyncResourceses().Delete(ctx, syncResource.Name, metav1.DeleteOptions{})
	assert.NoError(t, err)
	// 因为syncResource被syncRegistry所引用，所以调和失败，syncer不会被删除
	assert.Never(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return !ok
	}, 1*time.Second, 100*time.Millisecond)
}

func TestCreateTransformRule(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "test-sync-resource")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRule.TransformRefName = "test-transform-rule"
	syncResource := newSyncResource("test-sync-resource", []searchv1beta1.ResourceSyncRule{syncRule})
	transformRule := newTransformRule("test-transform-rule", "replace", "test-template")

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.syncResourcesLister = append(f.syncResourcesLister, syncResource)
	f.objects = append(f.objects, cluster, syncRegistry, syncResource)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.run(ctx, 1)

	_, err := f.client.SearchV1beta1().TransformRules().Create(ctx, transformRule, metav1.CreateOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	mgr, _ := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
	gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	assert.True(t, mgr.IsSyncerExist(gvr))
	sr := mgr.(*fakeSingleClusterSyncManager).GetSyncResources()
	assert.Equal(t, 1, len(sr))
	assert.Equal(t, *sr[gvr].Transform, transformRule.Spec)
}

func TestUpdateTransformRule(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "test-sync-resource")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRule.TransformRefName = "test-transform-rule"
	syncResource := newSyncResource("test-sync-resource", []searchv1beta1.ResourceSyncRule{syncRule})
	transformRule := newTransformRule("test-transform-rule", "replace", "test-template")

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.syncResourcesLister = append(f.syncResourcesLister, syncResource)
	f.transformRuleLister = append(f.transformRuleLister, transformRule)
	f.objects = append(f.objects, cluster, syncRegistry, syncResource, transformRule)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.run(ctx, 1)

	transformRule, err := f.client.SearchV1beta1().TransformRules().Get(ctx, transformRule.Name, metav1.GetOptions{})
	assert.NoError(t, err)
	transformRuleCopy := transformRule.DeepCopy()
	transformRuleCopy.Spec.Type = "patch"
	transformRuleCopy.Spec.ValueTemplate = "test-template2"
	transformRuleCopy.ResourceVersion = "1"
	_, err = f.client.SearchV1beta1().TransformRules().Update(ctx, transformRuleCopy, metav1.UpdateOptions{})
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		if mgr, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name); ok {
			gvr := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
			return mgr.IsSyncerExist(gvr) && *mgr.(*fakeSingleClusterSyncManager).GetSyncResources()[gvr].Transform == transformRuleCopy.Spec
		}
		return false

	}, wait.ForeverTestTimeout, 100*time.Millisecond)
}

func TestDeleteTransformRule(t *testing.T) {
	f := newFixture(t)
	cluster := newCluster("test-cluster", "app/test")
	insecure := true
	cluster.Spec.Access = clusterv1beta1.ClusterAccess{
		Endpoint: "https://test-cluster",
		Insecure: &insecure,
	}
	syncRegistry := newSyncRegistry("test-sync-registry", []string{"test-cluster"}, "test-sync-resource")
	syncRule := newResourceSyncRule("v1", "pods", "default")
	syncRule.TransformRefName = "test-transform-rule"
	syncResource := newSyncResource("test-sync-resource", []searchv1beta1.ResourceSyncRule{syncRule})
	transformRule := newTransformRule("test-transform-rule", "replace", "test-template")

	_, ctx := ktesting.NewTestContext(t)

	f.clusterLister = append(f.clusterLister, cluster)
	f.syncRegistryLister = append(f.syncRegistryLister, syncRegistry)
	f.syncResourcesLister = append(f.syncResourcesLister, syncResource)
	f.objects = append(f.objects, cluster, syncRegistry, syncResource)
	c, factory := f.newSyncController(ctx)
	factory.Start(ctx.Done())
	go c.run(ctx, 1)

	_, err := f.client.SearchV1beta1().TransformRules().Create(ctx, transformRule, metav1.CreateOptions{})
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return ok
	}, wait.ForeverTestTimeout, 100*time.Millisecond)

	err = f.client.SearchV1beta1().TransformRules().Delete(ctx, transformRule.Name, metav1.DeleteOptions{})
	assert.NoError(t, err)
	// 因为transformRule被syncResource所引用，所以调和失败，syncer不会被删除
	assert.Never(t, func() bool {
		_, ok := c.multiClusterSyncManager.GetSingleClusterSyncManager(cluster.Name)
		return !ok
	}, 1*time.Second, 100*time.Millisecond)
}

func TestGetSyncRegistries(t *testing.T) {
	f := newFixture(t)
	controller, factory := f.newSyncController(context.Background())

	cluster_informer := factory.Cluster().V1beta1().Clusters()
	cluster_indexer := cluster_informer.Informer().GetIndexer()

	registry_informer := factory.Search().V1beta1().SyncRegistries()
	registry_indexer := registry_informer.Informer().GetIndexer()

	cluster := &clusterv1beta1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"app": "test",
			},
		},
	}

	syncRegistries := []*searchv1beta1.SyncRegistry{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test1",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters: []string{"test"},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test2",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				ClusterLabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test3",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters: []string{"test1"},
				ClusterLabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test4",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters: []string{"test1"},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test5",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				ClusterLabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test1",
					},
				},
			},
		},
	}

	tests := []struct {
		name           string
		cluster        *clusterv1beta1.Cluster
		syncRegistries []*searchv1beta1.SyncRegistry
		expected       []*searchv1beta1.SyncRegistry
	}{
		{
			name:           "test1",
			cluster:        cluster,
			syncRegistries: syncRegistries,
			expected: []*searchv1beta1.SyncRegistry{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test1",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"test"},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test2",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						ClusterLabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "test",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test3",
					},
					Spec: searchv1beta1.SyncRegistrySpec{
						Clusters: []string{"test1"},
						ClusterLabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "test",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		cluster_indexer.Add(test.cluster)
		for _, syncRegistry := range test.syncRegistries {
			registry_indexer.Add(syncRegistry)
		}

		ret, err := controller.getSyncRegistries(test.cluster)
		assert.NoError(t, err)
		assert.ElementsMatch(t, test.expected, ret)
	}

}

func TestGetClusterSyncResources(t *testing.T) {
	f := newFixture(t)
	controller, factory := f.newSyncController(context.Background())

	cluster_informer := factory.Cluster().V1beta1().Clusters()
	cluster_indexer := cluster_informer.Informer().GetIndexer()

	registry_informer := factory.Search().V1beta1().SyncRegistries()
	registry_indexer := registry_informer.Informer().GetIndexer()

	resources_informer := factory.Search().V1beta1().SyncResourceses()
	resources_indexer := resources_informer.Informer().GetIndexer()

	cluster := &clusterv1beta1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
			Labels: map[string]string{
				"app": "test",
			},
		},
	}

	resourceSyncRule := map[string]searchv1beta1.ResourceSyncRule{
		"pods": {
			APIVersion: "v1",
			Resource:   "pods",
			Namespace:  "default",
		},
		"deployments": {
			APIVersion: "v1",
			Resource:   "deployments",
			Namespace:  "default_ns",
		},
		"services": {
			APIVersion: "v1",
			Resource:   "services",
			Namespace:  "default",
		},
	}

	syncRegistries := []*searchv1beta1.SyncRegistry{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters:             []string{"test-cluster"},
				SyncResourcesRefName: "test-sync-resources",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry-2",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				ClusterLabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				},
				SyncResources: []searchv1beta1.ResourceSyncRule{
					resourceSyncRule["deployments"],
				},
				SyncResourcesRefName: "test-sync-resources-2",
			},
		},
	}

	syncResources := []*searchv1beta1.SyncResources{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-resources",
			},
			Spec: searchv1beta1.SyncResourcesSpec{
				SyncResources: []searchv1beta1.ResourceSyncRule{
					resourceSyncRule["pods"],
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-resources-2",
			},
			Spec: searchv1beta1.SyncResourcesSpec{
				SyncResources: []searchv1beta1.ResourceSyncRule{
					resourceSyncRule["services"],
				},
			},
		},
	}

	tests := []struct {
		name           string
		cluster        *clusterv1beta1.Cluster
		syncRegistries []*searchv1beta1.SyncRegistry
		syncResources  []*searchv1beta1.SyncResources
		expected       map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule
	}{
		{
			name:           "test1",
			cluster:        cluster,
			syncRegistries: syncRegistries,
			syncResources:  syncResources,
			expected: map[schema.GroupVersionResource]*searchv1beta1.ResourceSyncRule{
				{
					Group:    "",
					Version:  "v1",
					Resource: "pods",
				}: {
					APIVersion: "v1",
					Resource:   "pods",
					Namespace:  "default",
				}, {
					Group:    "",
					Version:  "v1",
					Resource: "deployments",
				}: {
					APIVersion: "v1",
					Resource:   "deployments",
					Namespace:  "default_ns",
				}, {
					Group:    "",
					Version:  "v1",
					Resource: "services",
				}: {
					APIVersion: "v1",
					Resource:   "services",
					Namespace:  "default",
				},
			},
		},
	}

	for _, test := range tests {
		cluster_indexer.Add(test.cluster)
		for _, syncRegistry := range test.syncRegistries {
			registry_indexer.Add(syncRegistry)
		}
		for _, syncResources := range test.syncResources {
			resources_indexer.Add(syncResources)
		}
		ret, err := controller.getClusterSyncResources(test.cluster)
		assert.NoError(t, err)

		for k, v := range test.expected {
			assert.Equal(t, *v, *ret[k])
			delete(ret, k)
		}

		assert.Equal(t, 0, len(ret))
	}
}

func TestGetMatchedClusterNamesBySyncResources(t *testing.T) {
	f := newFixture(t)
	controller, factory := f.newSyncController(context.Background())

	cluster_informer := factory.Cluster().V1beta1().Clusters()
	cluster_indexer := cluster_informer.Informer().GetIndexer()

	registry_informer := factory.Search().V1beta1().SyncRegistries()
	registry_indexer := registry_informer.Informer().GetIndexer()

	resources_informer := factory.Search().V1beta1().SyncResourceses()
	resources_indexer := resources_informer.Informer().GetIndexer()

	clusters := []*clusterv1beta1.Cluster{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-2",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-3",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-4",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-5",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
	}

	syncRegistries := []*searchv1beta1.SyncRegistry{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters:             []string{"test-cluster"},
				SyncResourcesRefName: "test-sync-resources",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry-2",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				ClusterLabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				},
				SyncResourcesRefName: "test-sync-resources",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry-3",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters:             []string{"test-cluster-4"},
				SyncResourcesRefName: "test-sync-resources-2",
			},
		},
	}

	syncResources := &searchv1beta1.SyncResources{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-sync-resources",
		},
		Spec: searchv1beta1.SyncResourcesSpec{},
	}

	tests := []struct {
		name           string
		clusters       []*clusterv1beta1.Cluster
		syncRegistries []*searchv1beta1.SyncRegistry
		syncResources  *searchv1beta1.SyncResources
		expected       []string
	}{
		{
			name:           "test1",
			clusters:       clusters,
			syncRegistries: syncRegistries,
			syncResources:  syncResources,
			expected:       []string{"test-cluster", "test-cluster-2", "test-cluster-3", "test-cluster-5"},
		},
	}

	for _, test := range tests {
		for _, cluster := range test.clusters {
			cluster_indexer.Add(cluster)
		}
		for _, syncRegistry := range test.syncRegistries {
			registry_indexer.Add(syncRegistry)
		}
		resources_indexer.Add(test.syncResources)

		ret := controller.getMatchedClusterNamesBySyncResources(test.syncResources)
		// 判断两个切片中的元素是否相同，不考虑顺序
		assert.ElementsMatch(t, test.expected, ret)
	}
}

func TestGetMatchedClusterNamesByTransFormRule(t *testing.T) {
	f := newFixture(t)
	controller, factory := f.newSyncController(context.Background())

	cluster_informer := factory.Cluster().V1beta1().Clusters()
	cluster_indexer := cluster_informer.Informer().GetIndexer()

	registry_informer := factory.Search().V1beta1().SyncRegistries()
	registry_indexer := registry_informer.Informer().GetIndexer()

	resources_informer := factory.Search().V1beta1().SyncResourceses()
	resources_indexer := resources_informer.Informer().GetIndexer()

	transformRule_informer := factory.Search().V1beta1().TransformRules()
	transformRule_indexer := transformRule_informer.Informer().GetIndexer()

	clusters := []*clusterv1beta1.Cluster{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-2",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-3",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-4",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster-5",
				Labels: map[string]string{
					"app": "test",
				},
			},
		},
	}

	syncRegistries := []*searchv1beta1.SyncRegistry{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters:             []string{"test-cluster", "test-cluster-2", "test-cluster-3"},
				SyncResourcesRefName: "test-sync-resources",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry-2",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				ClusterLabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "test",
					},
				},
				SyncResourcesRefName: "test-sync-resources-2",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-registry-3",
			},
			Spec: searchv1beta1.SyncRegistrySpec{
				Clusters:             []string{"test-cluster-4"},
				SyncResourcesRefName: "test-sync-resources-3",
			},
		},
	}

	syncResources := []*searchv1beta1.SyncResources{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-resources",
			},
			Spec: searchv1beta1.SyncResourcesSpec{
				SyncResources: []searchv1beta1.ResourceSyncRule{
					{
						APIVersion:       "v1",
						Resource:         "pods",
						Namespace:        "default",
						TransformRefName: "test-transform-rule",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-resources-2",
			},
			Spec: searchv1beta1.SyncResourcesSpec{
				SyncResources: []searchv1beta1.ResourceSyncRule{
					{
						APIVersion:       "v1",
						Resource:         "deployments",
						Namespace:        "default",
						TransformRefName: "test-transform-rule",
					},
					{
						APIVersion:       "v1",
						Resource:         "services",
						Namespace:        "default",
						TransformRefName: "test-transform-rule-2",
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-sync-resources-3",
			},
			Spec: searchv1beta1.SyncResourcesSpec{
				SyncResources: []searchv1beta1.ResourceSyncRule{
					{
						APIVersion:       "v1",
						Resource:         "services",
						Namespace:        "default",
						TransformRefName: "test-transform-rule-2",
					},
					{
						APIVersion:       "v1",
						Resource:         "services",
						Namespace:        "default",
						TransformRefName: "test-transform-rule-3",
					},
				},
			},
		},
	}

	transformRule := &searchv1beta1.TransformRule{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-transform-rule",
		},
		Spec: searchv1beta1.TransformRuleSpec{},
	}

	tests := []struct {
		name                  string
		clusters              []*clusterv1beta1.Cluster
		syncRegistries        []*searchv1beta1.SyncRegistry
		syncResources         []*searchv1beta1.SyncResources
		transformRule         *searchv1beta1.TransformRule
		expectedSyncResources []*searchv1beta1.SyncResources
		expectedClusterNames  []string
	}{
		{
			name:           "test1",
			clusters:       clusters,
			syncRegistries: syncRegistries,
			syncResources:  syncResources,
			transformRule:  transformRule,
			expectedSyncResources: []*searchv1beta1.SyncResources{
				syncResources[0], syncResources[1],
			},
			expectedClusterNames: []string{"test-cluster", "test-cluster-2", "test-cluster-3", "test-cluster-5"},
		},
	}

	for _, test := range tests {
		for _, cluster := range test.clusters {
			cluster_indexer.Add(cluster)
		}
		for _, syncRegistry := range test.syncRegistries {
			registry_indexer.Add(syncRegistry)
		}
		for _, syncResource := range test.syncResources {
			resources_indexer.Add(syncResource)
		}
		if test.transformRule != nil {
			transformRule_indexer.Add(test.transformRule)
		}

		ret_resources := controller.getMatchedSyncResources(test.transformRule)
		assert.ElementsMatch(t, test.expectedSyncResources, ret_resources)

		ret := controller.getMatchedClusterNamesByTransFormRule(test.transformRule)
		assert.ElementsMatch(t, test.expectedClusterNames, ret)

	}

}
