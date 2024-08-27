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
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/kubectl/pkg/cmd/get"
)

type (
	Parser         = jsonpath.Parser
	Node           = jsonpath.Node
	ListNode       = jsonpath.ListNode
	TextNode       = jsonpath.TextNode
	FieldNode      = jsonpath.FieldNode
	ArrayNode      = jsonpath.ArrayNode
	FilterNode     = jsonpath.FilterNode
	IntNode        = jsonpath.IntNode
	BoolNode       = jsonpath.BoolNode
	FloatNode      = jsonpath.FloatNode
	WildcardNode   = jsonpath.WildcardNode
	RecursiveNode  = jsonpath.RecursiveNode
	UnionNode      = jsonpath.UnionNode
	IdentifierNode = jsonpath.IdentifierNode
)

var (
	Parse                     = jsonpath.Parse
	RelaxedJSONPathExpression = get.RelaxedJSONPathExpression
)
