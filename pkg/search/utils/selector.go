package utils

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

type Selectable interface {
	GetLabels() labels.Labels
	GetFields() fields.Fields
}

type Selector struct {
	Label labels.Selector
	Field *FieldsSelector
}

func (s Selector) ServerSupported() bool {
	if s.Field != nil && !s.Field.ServerSupported {
		return false
	}
	return true
}

type FieldsSelector struct {
	fields.Selector
	ServerSupported bool
}

func (s Selector) Matches(obj Selectable) (bool, error) {
	if s.Label != nil {
		if !s.Label.Matches(obj.GetLabels()) {
			return false, nil
		}
	}
	if s.Field != nil {
		if !s.Field.Matches(obj.GetFields()) {
			return false, nil
		}
	}
	return true, nil
}

type MultiSelectors []Selector

func (m MultiSelectors) Matches(obj Selectable) (bool, error) {
	for _, s := range m {
		matches, err := s.Matches(obj)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

type selectableUnstructured struct {
	*unstructured.Unstructured
	parser *JSONPathParser
}

func (w selectableUnstructured) GetLabels() labels.Labels {
	return labels.Set(w.Unstructured.GetLabels())
}

func (w selectableUnstructured) GetFields() fields.Fields {
	return NewJSONPathFields(w.parser, w.UnstructuredContent())
}

func SelectableUnstructured(u *unstructured.Unstructured, parser *JSONPathParser) Selectable {
	return selectableUnstructured{Unstructured: u, parser: parser}
}
