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
	"fmt"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clientgocache "k8s.io/client-go/tools/cache"

	"github.com/KusionStack/karpor/pkg/kubernetes/apis/search/v1beta1"
	syncercache "github.com/KusionStack/karpor/pkg/syncer/cache"
	"github.com/KusionStack/karpor/pkg/syncer/jsonextracter"
	"github.com/KusionStack/karpor/pkg/syncer/transform"
	"github.com/KusionStack/karpor/pkg/util/jsonpath"
)

// parseTrimer creates and returns a trim function for the informerSource based on the provider trims.
func parseTrimer(ctx context.Context, t *v1beta1.TrimRuleSpec) (syncercache.TransformFunc, error) {
	if t == nil || len(t.Retain.JSONPaths) == 0 {
		//nolint:nilnil,nolintlint
		return nil, nil
	}

	extracters := make([]jsonextracter.Extracter, 0, len(t.Retain.JSONPaths))
	for _, p := range t.Retain.JSONPaths {
		p, err := jsonpath.RelaxedJSONPathExpression(p)
		if err != nil {
			return nil, err
		}

		ex, err := jsonextracter.BuildExtracter(p, true)
		if err != nil {
			return nil, err
		}
		extracters = append(extracters, ex)
	}

	trimFunc := func(obj interface{}) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				logr.FromContext(ctx).Error(err, "error in triming object")
				ret, err = obj, nil
			}
		}()

		if d, ok := obj.(clientgocache.DeletedFinalStateUnknown); ok {
			// Since we import ES data into informer cache at startup, the
			// resource that was deleted during the restart will generate
			// DeletedFinalStateUnknown.
			// We unwarp the object here, so there is no need for following
			// steps including event handler to care about DeletedFinalStateUnknown.
			obj = d.Obj
		}

		u, ok := obj.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("trim: object's type should be *unstructured.Unstructured, but received %T", obj)
		}

		merged, err := jsonextracter.Merge(extracters, u.Object)
		if err != nil {
			return nil, err
		}

		unObj := &unstructured.Unstructured{Object: merged}
		return unObj, nil
	}

	return trimFunc, nil
}

// parseTransformer creates and returns a transformation function for the informerSource based on the provider transformers.
func parseTransformer(ctx context.Context, t *v1beta1.TransformRuleSpec, clusterName string) (syncercache.TransformFunc, error) {
	if t == nil {
		//nolint:nilnil,nolintlint
		return nil, nil
	}

	fn, found := transform.GetTransformFunc(t.Type)
	if !found {
		return nil, fmt.Errorf("unsupported transform type %q", t.Type)
	}

	tmpl, err := newTemplate(t.ValueTemplate, clusterName)
	if err != nil {
		return nil, errors.Wrap(err, "invalid transform template")
	}

	return func(obj interface{}) (ret interface{}, err error) {
		defer func() {
			if err != nil {
				logr.FromContext(ctx).Error(err, "error in transforming object")
			}
		}()

		u, ok := obj.(*unstructured.Unstructured)
		if !ok {
			return nil, fmt.Errorf("transform: object's type should be *unstructured.Unstructured, but received %T", obj)
		}

		templateData := struct {
			*unstructured.Unstructured
			Cluster string
		}{
			Unstructured: u,
			Cluster:      clusterName,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, templateData); err != nil {
			return nil, errors.Wrap(err, "transform: error rendering template")
		}
		return fn(obj, buf.String())
	}, nil
}

// newTemplate creates and returns a new text template from the provided string, which can be used for processing templates in the syncer.
func newTemplate(tmpl, cluster string) (*template.Template, error) {
	clusterFuncs, _ := transform.GetClusterTmplFuncs(cluster)
	return template.New("transformTemplate").Funcs(sprig.FuncMap()).Funcs(clusterFuncs).Parse(tmpl)
}
