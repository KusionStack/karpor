package syncer

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
)

type AgentReconciler struct {
	SyncReconciler

	clusterName string
}

// NewAgentReconciler creates a new instance of the AgentReconciler structure with the given storage.
func NewAgentReconciler(storage storage.ResourceStorage, highAvailability bool, clusterName string) *AgentReconciler {
	return &AgentReconciler{
		SyncReconciler: SyncReconciler{
			storage:          storage,
			highAvailability: highAvailability,
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
		return reconcile.Result{}, r.stopCluster(ctx, cluster.Name)
	}

	return reconcile.Result{}, r.handleClusterAddOrUpdate(ctx, cluster.DeepCopy())
}
