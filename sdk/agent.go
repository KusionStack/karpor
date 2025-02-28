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

package sdk

import (
    "context"

    esclient "github.com/elastic/go-elasticsearch/v8"
    "github.com/pkg/errors"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/dynamic"
    "k8s.io/klog/v2/klogr"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/manager"

    "github.com/KusionStack/karpor/pkg/infra/search/storage/elasticsearch"
    "github.com/KusionStack/karpor/pkg/syncer"
    "github.com/KusionStack/karpor/pkg/syncer/utils"
)

type AgentOptions struct {
    ElasticSearchAddresses []string
    ClusterName            string

    SyncGVKs []schema.GroupVersionKind
}

/*
StartAgentWithOperator defines sdk entrance.
Developers can use this function to share memory with exited operator for special resources in SyncGVK.
Only support push mode.
*/
func StartAgentWithOperator(ctx context.Context, existMgr manager.Manager, options *AgentOptions) error {
    ctrl.SetLogger(klogr.New())
    log := ctrl.Log.WithName("setup")

    defer func() {
        if err := recover(); err != nil {
            log.Error(nil, "recovered panic in karpor agent", err)
        }
    }()

    // apply crds
    dynamicClient, err := dynamic.NewForConfig(ctrl.GetConfigOrDie())
    if err != nil {
        return errors.Wrapf(err, "failed to build dynamic client for ageng")
    }
    err = utils.ApplyCrds(ctx, dynamicClient)
    if err != nil {
        return err
    }

    // TODO: add startup parameters to change the type of storage
    //nolint:contextcheck
    es, err := elasticsearch.NewStorage(esclient.Config{
        Addresses: options.ElasticSearchAddresses,
    })
    if err != nil {
        log.Error(err, "unable to init elasticsearch client")
        return err
    }

    // start operator for special resource
    if len(options.SyncGVKs) != 0 {
        for idx := range options.SyncGVKs {
            gvk := options.SyncGVKs[idx]
            utils.SetSyncGVK(gvk, utils.ZeroVal)
            err := syncer.NewDynamicReconciler(ctx, options.ClusterName, gvk, es).SetupWithManager(existMgr)
            if err != nil {
                return err
            }
        }
    }

    //nolint:contextcheck
    if err = syncer.NewAgentReconciler(es, options.ClusterName).SetupWithManager(existMgr); err != nil {
        log.Error(err, "unable to create resource syncer")
        return err
    }

    return nil
}
