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

package exmaple

import (
    "context"
    "fmt"
    "os"

    "k8s.io/apimachinery/pkg/runtime/schema"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"

    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"

    "github.com/KusionStack/karpor/pkg/kubernetes/scheme"
    "github.com/KusionStack/karpor/sdk"
)

type PodReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    pod := &corev1.Pod{}
    if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    fmt.Printf("Reconciling Pod: %s/%s\n", pod.Namespace, pod.Name)

    return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&corev1.Pod{}).
        Complete(r)
}

// example for sdk, only support push mode.
func main() {
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme: scheme.Scheme,
    })
    if err != nil {
        fmt.Fprintf(os.Stderr, "unable to start manager: %v\n", err)
        os.Exit(1)
    }

    if err = (&PodReconciler{
        Client: mgr.GetClient(),
        Scheme: mgr.GetScheme(),
    }).SetupWithManager(mgr); err != nil {
        fmt.Fprintf(os.Stderr, "unable to create controller: %v\n", err)
        os.Exit(1)
    }

    err = sdk.StartAgentWithOperator(context.Background(), mgr, &sdk.AgentOptions{
        ClusterName:            "example-cluster",
        ElasticSearchAddresses: []string{"http://127.0.0.1:9200"},
        SyncGVKs: []schema.GroupVersionKind{
            {Group: "", Version: "v1", Kind: "Pod"},
        },
    })
    if err != nil {
        fmt.Fprintf(os.Stderr, "unable to add dynamic controller: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Starting manager")
    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        fmt.Fprintf(os.Stderr, "problem running manager: %v\n", err)
        os.Exit(1)
    }
}
