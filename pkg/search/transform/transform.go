package transform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	Register("patch", Patch)
	Register("repalce", Replace)
}

type Transformer struct {
	fn      TransformFunc
	t       *template.Template
	cluster string
}

func NewTransformer(tType string, tmpl string, cluster string) (*Transformer, error) {
	fn, found := GetTransformFunc(tType)
	if !found {
		return nil, fmt.Errorf("unsupported transform type %q", tType)
	}
	t, err := NewTemplate(tmpl)
	if err != nil {
		return nil, err
	}
	return &Transformer{fn: fn, t: t, cluster: cluster}, nil
}

func (t *Transformer) Transform(original interface{}) (interface{}, error) {
	var buf bytes.Buffer
	if err := t.t.Execute(&buf, templateData{Obj: original, Cluster: t.cluster}); err != nil {
		return nil, fmt.Errorf("error rendering template: %v", err)
	}
	return t.fn(original, buf.String())
}

type templateData struct {
	Obj     interface{} `json:"obj"`
	Cluster string      `json:"cluster"`
}

func NewTemplate(tmpl string) (*template.Template, error) {
	return template.New("transformTemplate").Funcs(sprig.FuncMap()).Parse(tmpl)
}

var defaultRegistry TransformFuncRegistry

func Register(tType string, transFunc TransformFunc) {
	defaultRegistry.Register(tType, transFunc)
}

func GetTransformFunc(transformerType string) (TransformFunc, bool) {
	return defaultRegistry.Get(transformerType)
}

type TransformFuncRegistry struct {
	transformers map[string]TransformFunc
}

func (r *TransformFuncRegistry) Register(tType string, transFunc TransformFunc) {
	r.transformers[tType] = transFunc
}

func (r *TransformFuncRegistry) Get(transformerType string) (transFunc TransformFunc, found bool) {
	transFunc, found = r.transformers[transformerType]
	return
}

type TransformFunc func(original interface{}, data string) (target interface{}, err error)

func Patch(original interface{}, patchText string) (interface{}, error) {
	patch, err := jsonpatch.DecodePatch([]byte(patchText))
	if err != nil {
		return nil, fmt.Errorf("patch is invalid: %v", err)
	}

	originalJSON, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	modifiedJSON, err := patch.Apply(originalJSON)
	if err != nil {
		return nil, err
	}

	var dest interface{}
	u, ok := original.(runtime.Unstructured)
	if !ok {
		return nil, fmt.Errorf(`type "%vT not supported`, original)
	}
	dest = u.NewEmptyInstance()

	if err := json.Unmarshal(modifiedJSON, &dest); err != nil {
		return nil, fmt.Errorf("json decoding error: %v", err)
	}
	return &dest, nil
}

func Replace(original interface{}, jsonString string) (interface{}, error) {
	var dest unstructured.Unstructured
	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		return nil, fmt.Errorf("json decoding error: %v", err)
	}
	return &dest, nil
}
