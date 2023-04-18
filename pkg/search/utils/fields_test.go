package utils

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestJSONPathFields(t *testing.T) {
	u := unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": metav1.ObjectMeta{
				Name: "foo",
				Labels: map[string]string{
					"label1":        "bar",
					"label2/label3": "foo",
				},
			},
		},
	}

	fields := NewJSONPathFields(NewJSONPathParser(), u.Object)

	testData := map[string]string{
		`{.metadata.name}`:                    "foo",
		`metadata.name`:                       "foo",
		`.metadata.name`:                      "foo",
		`{.metadata.labels.label1}`:           "bar",
		`{.metadata.labels['label1']}`:        "bar",
		`{.metadata.labels['label2/label3']}`: "foo",
		`{.notExistField}`:                    "",
	}

	for path, expectVal := range testData {
		actualVal := fields.Get(path)
		if actualVal != expectVal {
			t.Errorf(`the value of path '%s' is expected to be '%s', but got '%s'`, path, expectVal, actualVal)
		}
	}
}
