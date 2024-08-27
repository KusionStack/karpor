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

package jsonextracter

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/client-go/third_party/forked/golang/template"
)

type JSONPath struct {
	name       string
	parser     *Parser
	beginRange int
	inRange    int
	endRange   int

	lastEndNode *Node

	allowMissingKeys bool
}

// New creates a new JSONPath with the given name.
func New(name string) *JSONPath {
	return &JSONPath{
		name:       name,
		beginRange: 0,
		inRange:    0,
		endRange:   0,
	}
}

// AllowMissingKeys allows a caller to specify whether they want an error if a field or map key
// cannot be located, or simply an empty result. The receiver is returned for chaining.
func (j *JSONPath) AllowMissingKeys(allow bool) *JSONPath {
	j.allowMissingKeys = allow
	return j
}

// Parse parses the given template and returns an error.
func (j *JSONPath) Parse(text string) error {
	var err error
	j.parser, err = Parse(j.name, text)
	return err
}

type setFieldFunc func(val reflect.Value) error

var nopSetFieldFunc = func(_ reflect.Value) error { return nil }

func makeNopSetFieldFuncSlice(n int) []setFieldFunc {
	fns := make([]setFieldFunc, n)
	for i := 0; i < n; i++ {
		fns[i] = nopSetFieldFunc
	}
	return fns
}

// Extract outputs the field specified by JSONPath.
// The output contains not only the field value, but also its upstream structure.
//
// The data structure of the extracted field must be of type `map[string]interface{}`,
// and `struct` is not supported (an error will be returned).
func (j *JSONPath) Extract(data map[string]interface{}) (map[string]interface{}, error) {
	container := struct{ Root reflect.Value }{}
	setFn := func(val reflect.Value) error {
		container.Root = val
		return nil
	}

	_, err := j.FindResults(data, setFn)
	if err != nil {
		return nil, err
	}

	if !container.Root.IsValid() {
		return nil, nil
	}

	return container.Root.Interface().(map[string]interface{}), nil
}

func (j *JSONPath) FindResults(data interface{}, setFn setFieldFunc) ([][]reflect.Value, error) {
	if j.parser == nil {
		return nil, fmt.Errorf("%s is an incomplete jsonpath template", j.name)
	}

	cur := []reflect.Value{reflect.ValueOf(data)}
	curnFn := []setFieldFunc{setFn}
	nodes := j.parser.Root.Nodes
	fullResult := [][]reflect.Value{}
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		results, fn, err := j._walk(cur, node, curnFn)
		if err != nil {
			return nil, err
		}

		// encounter an end node, break the current block
		if j.endRange > 0 && j.endRange <= j.inRange {
			j.endRange--
			j.lastEndNode = &nodes[i]
			break
		}
		// encounter a range node, start a range loop
		if j.beginRange > 0 {
			j.beginRange--
			j.inRange++
			if len(results) > 0 {
				for ri, value := range results {
					j.parser.Root.Nodes = nodes[i+1:]
					nextResults, err := j.FindResults(value.Interface(), fn[ri])
					if err != nil {
						return nil, err
					}
					fullResult = append(fullResult, nextResults...)
				}
			} else {
				// If the range has no results, we still need to process the nodes within the range
				// so the position will advance to the end node
				j.parser.Root.Nodes = nodes[i+1:]
				_, err := j.FindResults(nil, nopSetFieldFunc)
				if err != nil {
					return nil, err
				}
			}
			j.inRange--

			// Fast forward to resume processing after the most recent end node that was encountered
			for k := i + 1; k < len(nodes); k++ {
				if &nodes[k] == j.lastEndNode {
					i = k
					break
				}
			}
			continue
		}
		fullResult = append(fullResult, results)
	}
	return fullResult, nil
}

func (j *JSONPath) _walk(value []reflect.Value, node Node, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	switch node := node.(type) {
	case *ListNode:
		return j._evalList(value, node, setFn)
	case *FieldNode:
		return j.evalField(value, node, setFn)
	case *ArrayNode:
		return j.evalArray(value, node, setFn)
	case *IdentifierNode:
		return j.evalIdentifier(value, node, setFn)
	case *UnionNode:
		return j._evalUnion(value, node, setFn)
	case *FilterNode:
		return j.evalFilter(value, node, setFn)
	default:
		return nil, nil, fmt.Errorf("Extract does not support node %v", node)
	}
}

// walk visits tree rooted at the given node in DFS order
func (j *JSONPath) walk(value []reflect.Value, node Node) ([]reflect.Value, error) {
	switch node := node.(type) {
	case *ListNode:
		return j.evalList(value, node)
	case *TextNode:
		return []reflect.Value{reflect.ValueOf(node.Text)}, nil
	case *FieldNode:
		value, _, err := j.evalField(value, node, makeNopSetFieldFuncSlice(len(value)))
		return value, err
	case *ArrayNode:
		value, _, err := j.evalArray(value, node, makeNopSetFieldFuncSlice(len(value)))
		return value, err
	case *FilterNode:
		value, _, err := j.evalFilter(value, node, makeNopSetFieldFuncSlice(len(value)))
		return value, err
	case *IntNode:
		return j.evalInt(value, node)
	case *BoolNode:
		return j.evalBool(value, node)
	case *FloatNode:
		return j.evalFloat(value, node)
	case *WildcardNode:
		return j.evalWildcard(value, node)
	case *RecursiveNode:
		return j.evalRecursive(value, node)
	case *UnionNode:
		return j.evalUnion(value, node)
	case *IdentifierNode:
		value, _, err := j.evalIdentifier(value, node, makeNopSetFieldFuncSlice(len(value)))
		return value, err
	default:
		return value, fmt.Errorf("unexpected Node %v", node)
	}
}

// evalInt evaluates IntNode
func (j *JSONPath) evalInt(input []reflect.Value, node *IntNode) ([]reflect.Value, error) {
	result := make([]reflect.Value, len(input))
	for i := range input {
		result[i] = reflect.ValueOf(node.Value)
	}
	return result, nil
}

// evalFloat evaluates FloatNode
func (j *JSONPath) evalFloat(input []reflect.Value, node *FloatNode) ([]reflect.Value, error) {
	result := make([]reflect.Value, len(input))
	for i := range input {
		result[i] = reflect.ValueOf(node.Value)
	}
	return result, nil
}

// evalBool evaluates BoolNode
func (j *JSONPath) evalBool(input []reflect.Value, node *BoolNode) ([]reflect.Value, error) {
	result := make([]reflect.Value, len(input))
	for i := range input {
		result[i] = reflect.ValueOf(node.Value)
	}
	return result, nil
}

func (j *JSONPath) _evalList(value []reflect.Value, node *ListNode, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	var err error
	curValue := value
	curFns := setFn

	for _, node := range node.Nodes {
		curValue, curFns, err = j._walk(curValue, node, curFns)
		if err != nil {
			return curValue, curFns, err
		}
	}
	return curValue, curFns, nil
}

// evalList evaluates ListNode
func (j *JSONPath) evalList(value []reflect.Value, node *ListNode) ([]reflect.Value, error) {
	var err error
	curValue := value
	for _, node := range node.Nodes {
		curValue, err = j.walk(curValue, node)
		if err != nil {
			return curValue, err
		}
	}
	return curValue, nil
}

// evalIdentifier evaluates IdentifierNode
func (j *JSONPath) evalIdentifier(input []reflect.Value, node *IdentifierNode, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	results := []reflect.Value{}
	switch node.Name {
	case "range":
		j.beginRange++
		results = input
	case "end":
		if j.inRange > 0 {
			j.endRange++
		} else {
			return results, setFn, fmt.Errorf("not in range, nothing to end")
		}
	default:
		return input, setFn, fmt.Errorf("unrecognized identifier %v", node.Name)
	}
	return results, setFn, nil
}

// evalArray evaluates ArrayNode
func (j *JSONPath) evalArray(input []reflect.Value, node *ArrayNode, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	result := []reflect.Value{}
	nextFns := []setFieldFunc{}
	for k, value := range input {

		value, isNil := template.Indirect(value)
		if isNil {
			continue
		}
		if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
			return input, nextFns, fmt.Errorf("%v is not array or slice", value.Type())
		}
		params := node.Params
		if !params[0].Known {
			params[0].Value = 0
		}
		if params[0].Value < 0 {
			params[0].Value += value.Len()
		}
		if !params[1].Known {
			params[1].Value = value.Len()
		}

		if params[1].Value < 0 || (params[1].Value == 0 && params[1].Derived) {
			params[1].Value += value.Len()
		}
		sliceLength := value.Len()
		if params[1].Value != params[0].Value { // if you're requesting zero elements, allow it through.
			if params[0].Value >= sliceLength || params[0].Value < 0 {
				return input, nextFns, fmt.Errorf("array index out of bounds: index %d, length %d", params[0].Value, sliceLength)
			}
			if params[1].Value > sliceLength || params[1].Value < 0 {
				return input, nextFns, fmt.Errorf("array index out of bounds: index %d, length %d", params[1].Value-1, sliceLength)
			}
			if params[0].Value > params[1].Value {
				return input, nextFns, fmt.Errorf("starting index %d is greater than ending index %d", params[0].Value, params[1].Value)
			}
		} else {
			return result, nextFns, nil
		}

		value = value.Slice(params[0].Value, params[1].Value)

		step := 1
		if params[2].Known {
			if params[2].Value <= 0 {
				return input, nextFns, fmt.Errorf("step must be > 0")
			}
			step = params[2].Value
		}

		loopResult := []reflect.Value{}
		for i := 0; i < value.Len(); i += step {
			loopResult = append(loopResult, value.Index(i))
		}
		result = append(result, loopResult...)

		s := reflect.MakeSlice(value.Type(), len(loopResult), len(loopResult))
		for i := 0; i < len(loopResult); i++ {
			ii := i
			s.Index(ii).Set(loopResult[i])
			nextFns = append(nextFns, func(val reflect.Value) error {
				s.Index(ii).Set(val)
				return nil
			})
		}

		setFn[k](s)
	}
	return result, nextFns, nil
}

func (j *JSONPath) _evalUnion(input []reflect.Value, node *UnionNode, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	result := []reflect.Value{}
	fns := []setFieldFunc{}

	union := make([][]reflect.Value, len(input))
	setFn_ := make([]setFieldFunc, len(input))

	for i := 0; i < len(input); i++ {
		ii := i
		setFn_[i] = func(val reflect.Value) error {
			union[ii] = append(union[ii], val)
			return nil
		}
	}

	for _, listNode := range node.Nodes {
		temp, nextFn, err := j._evalList(input, listNode, setFn_)
		if err != nil {
			return input, fns, err
		}
		result = append(result, temp...)
		fns = append(fns, nextFn...)
	}

	for i, fn := range setFn {
		if len(union[i]) == 0 {
			continue
		}

		m := union[i][0]
		for j := 1; j < len(union[i]); j++ {
			val := union[i][j]
			for _, key := range val.MapKeys() {
				m.SetMapIndex(key, val.MapIndex(key))
			}
		}
		fn(m)
	}

	return result, fns, nil
}

// evalUnion evaluates UnionNode
func (j *JSONPath) evalUnion(input []reflect.Value, node *UnionNode) ([]reflect.Value, error) {
	result := []reflect.Value{}
	for _, listNode := range node.Nodes {
		temp, err := j.evalList(input, listNode)
		if err != nil {
			return input, err
		}
		result = append(result, temp...)
	}
	return result, nil
}

//lint:ignore U1000 ignore unused function
func (j *JSONPath) findFieldInValue(value *reflect.Value, node *FieldNode) (reflect.Value, error) {
	t := value.Type()
	var inlineValue *reflect.Value
	for ix := 0; ix < t.NumField(); ix++ {
		f := t.Field(ix)
		jsonTag := f.Tag.Get("json")
		parts := strings.Split(jsonTag, ",")
		if len(parts) == 0 {
			continue
		}
		if parts[0] == node.Value {
			return value.Field(ix), nil
		}
		if len(parts[0]) == 0 {
			val := value.Field(ix)
			inlineValue = &val
		}
	}
	if inlineValue != nil {
		if inlineValue.Kind() == reflect.Struct {
			// handle 'inline'
			match, err := j.findFieldInValue(inlineValue, node)
			if err != nil {
				return reflect.Value{}, err
			}
			if match.IsValid() {
				return match, nil
			}
		}
	}
	return value.FieldByName(node.Value), nil
}

// evalField evaluates field of struct or key of map.
func (j *JSONPath) evalField(input []reflect.Value, node *FieldNode, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	results := []reflect.Value{}
	nextFns := []setFieldFunc{}
	// If there's no input, there's no output
	if len(input) == 0 {
		return results, nextFns, nil
	}
	for k, value := range input {
		var result reflect.Value
		var fn setFieldFunc
		value, isNil := template.Indirect(value)
		if isNil {
			continue
		}

		if value.Kind() != reflect.Map {
			return results, nextFns, fmt.Errorf("%v is of the type %T, expected map[string]interface{}", value.Interface(), value.Interface())
		} else {
			mapKeyType := value.Type().Key()
			nodeValue := reflect.ValueOf(node.Value)
			// node value type must be convertible to map key type
			if !nodeValue.Type().ConvertibleTo(mapKeyType) {
				return results, nextFns, fmt.Errorf("%s is not convertible to %s", nodeValue, mapKeyType)
			}
			key := nodeValue.Convert(mapKeyType)
			result = value.MapIndex(key)

			val := reflect.MakeMap(value.Type())
			val.SetMapIndex(key, result)
			setFn[k](val)

			fn = func(val_ reflect.Value) error {
				val.SetMapIndex(key, val_)
				return nil
			}
		}

		if result.IsValid() {
			results = append(results, result)
			nextFns = append(nextFns, fn)
		}
	}
	if len(results) == 0 {
		if j.allowMissingKeys {
			return results, nextFns, nil
		}
		return results, nextFns, fmt.Errorf("%s is not found", node.Value)
	}
	return results, nextFns, nil
}

// evalWildcard extracts all contents of the given value
func (j *JSONPath) evalWildcard(input []reflect.Value, _ *WildcardNode) ([]reflect.Value, error) {
	results := []reflect.Value{}
	for _, value := range input {
		value, isNil := template.Indirect(value)
		if isNil {
			continue
		}

		kind := value.Kind()
		if kind == reflect.Struct {
			for i := 0; i < value.NumField(); i++ {
				results = append(results, value.Field(i))
			}
		} else if kind == reflect.Map {
			for _, key := range value.MapKeys() {
				results = append(results, value.MapIndex(key))
			}
		} else if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String {
			for i := 0; i < value.Len(); i++ {
				results = append(results, value.Index(i))
			}
		}
	}
	return results, nil
}

// evalRecursive visits the given value recursively and pushes all of them to result
func (j *JSONPath) evalRecursive(input []reflect.Value, node *RecursiveNode) ([]reflect.Value, error) {
	result := []reflect.Value{}
	for _, value := range input {
		results := []reflect.Value{}
		value, isNil := template.Indirect(value)
		if isNil {
			continue
		}

		kind := value.Kind()
		if kind == reflect.Struct {
			for i := 0; i < value.NumField(); i++ {
				results = append(results, value.Field(i))
			}
		} else if kind == reflect.Map {
			for _, key := range value.MapKeys() {
				results = append(results, value.MapIndex(key))
			}
		} else if kind == reflect.Array || kind == reflect.Slice || kind == reflect.String {
			for i := 0; i < value.Len(); i++ {
				results = append(results, value.Index(i))
			}
		}
		if len(results) != 0 {
			result = append(result, value)
			output, err := j.evalRecursive(results, node)
			if err != nil {
				return result, err
			}
			result = append(result, output...)
		}
	}
	return result, nil
}

// evalFilter filters array according to FilterNode
func (j *JSONPath) evalFilter(input []reflect.Value, node *FilterNode, setFn []setFieldFunc) ([]reflect.Value, []setFieldFunc, error) {
	results := []reflect.Value{}
	fns := []setFieldFunc{}
	for k, value := range input {
		value, _ = template.Indirect(value)

		if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
			return input, fns, fmt.Errorf("%v is not array or slice and cannot be filtered", value)
		}

		loopResult := []reflect.Value{}
		for i := 0; i < value.Len(); i++ {
			temp := []reflect.Value{value.Index(i)}
			lefts, err := j.evalList(temp, node.Left)

			// case exists
			if node.Operator == "exists" {
				if len(lefts) > 0 {
					results = append(results, value.Index(i))
				}
				continue
			}

			if err != nil {
				return input, fns, err
			}

			var left, right interface{}
			switch {
			case len(lefts) == 0:
				continue
			case len(lefts) > 1:
				return input, fns, fmt.Errorf("can only compare one element at a time")
			}
			left = lefts[0].Interface()

			rights, err := j.evalList(temp, node.Right)
			if err != nil {
				return input, fns, err
			}
			switch {
			case len(rights) == 0:
				continue
			case len(rights) > 1:
				return input, fns, fmt.Errorf("can only compare one element at a time")
			}
			right = rights[0].Interface()

			pass := false
			switch node.Operator {
			case "<":
				pass, err = template.Less(left, right)
			case ">":
				pass, err = template.Greater(left, right)
			case "==":
				pass, err = template.Equal(left, right)
			case "!=":
				pass, err = template.NotEqual(left, right)
			case "<=":
				pass, err = template.LessEqual(left, right)
			case ">=":
				pass, err = template.GreaterEqual(left, right)
			default:
				return results, fns, fmt.Errorf("unrecognized filter operator %s", node.Operator)
			}
			if err != nil {
				return results, fns, err
			}
			if pass {
				loopResult = append(loopResult, value.Index(i))
			}
		}

		s := reflect.MakeSlice(value.Type(), len(loopResult), len(loopResult))
		for i := 0; i < len(loopResult); i++ {
			ii := i
			s.Index(ii).Set(loopResult[i])
			fns = append(fns, func(val reflect.Value) error {
				s.Index(ii).Set(val)
				return nil
			})
		}

		setFn[k](s)
		results = append(results, loopResult...)
	}
	return results, fns, nil
}
