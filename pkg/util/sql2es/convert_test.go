// Copyright The Karbour Authors.
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

package sql2es

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xwb1989/sqlparser"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		sql           string
		expectedDSL   string
		expectedTable string
		expectedErr   string
	}{
		// Test cases with valid SQL queries
		{
			sql:           "SELECT * FROM mock_table WHERE id = '123'",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"match_phrase" : {"id" : {"query" : "123"}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name = 'John' AND age > 30",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"match_phrase" : {"name" : {"query" : "John"}}},{"range" : {"age" : {"gt" : "30"}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name = 'John' OR age > 30",
			expectedDSL:   `{"query" : {"bool" : {"should" : [{"match_phrase" : {"name" : {"query" : "John"}}},{"range" : {"age" : {"gt" : "30"}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name = ''",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"match_phrase" : {"name" : {"query" : ""}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE age >= 18 AND age <= 60",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"range" : {"age" : {"from" : "18"}}},{"range" : {"age" : {"to" : "60"}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name = 'John' AND age > 30 ORDER BY age DESC",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"match_phrase" : {"name" : {"query" : "John"}}},{"range" : {"age" : {"gt" : "30"}}}]}},"from" : 0,"size" : 1,"sort" : [{"age": "desc"}]}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name != 'abc'",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"bool" : {"must_not" : [{"match_phrase" : {"name" : {"query" : "abc"}}}]}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name in ('abc', 'def')",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"terms" : {"name" : ["abc", "def"]}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name not in ('abc', 'def')",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"bool" : {"must_not" : {"terms" : {"name" : ["abc", "def"]}}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE contains(name, 'abc')",
			expectedDSL:   `{"query" : {"match_phrase": {"name": "abc"}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name like '%abc%'",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"wildcard" : {"name" : "*abc*"}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name not like '%abc%'",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"bool" : {"must_not" : {"wildcard" : {"name" : "*abc*"}}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT count(*) FROM mock_table WHERE age > 18 AND age < 60",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"range" : {"age" : {"gt" : "18"}}},{"range" : {"age" : {"lt" : "60"}}}]}},"from" : 0,"size" : 0,"aggregations" : {"COUNT(*)":{"value_count":{"field":"_index"}}}}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT count(id) FROM mock_table WHERE age >= 18 AND age <= 60",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"range" : {"age" : {"from" : "18"}}},{"range" : {"age" : {"to" : "60"}}}]}},"from" : 0,"size" : 0,"aggregations" : {"COUNT(id)":{"value_count":{"field":"id"}}}}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE age between 18 and 60",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"range" : {"age" : {"from" : "18", "to" : "60"}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE age not between 18 and 60",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"range" : {"age" : {"from" : "18", "to" : "60"}}}]}},"from" : 0,"size" : 1}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name like '%abc%' GROUP BY date_histogram(field='create_time', value='1h')",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"wildcard" : {"name" : "*abc*"}}]}},"from" : 0,"size" : 0,"aggregations" : {"date_histogram(field=create_time,value=1h)":{"date_histogram":{"field":"create_time","format":"yyyy-MM-dd HH:mm:ss","interval":"1h"}}}}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name like '%abc%' GROUP BY range(age, 18, 60)",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"wildcard" : {"name" : "*abc*"}}]}},"from" : 0,"size" : 0,"aggregations" : {"range(age,18,60)":{"range":{"field":"age","ranges":[{"from":"18","to":"60"}]}}}}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name like '%abc%' GROUP BY date_range(field='create_time' , format='yyyy-MM-dd', 'now-8d','now-7d','now-6d','now')",
			expectedDSL:   `{"query" : {"bool" : {"must" : [{"wildcard" : {"name" : "*abc*"}}]}},"from" : 0,"size" : 0,"aggregations" : {"date_range(field=create_time,format=yyyy-MM-dd,now-8d,now-7d,now-6d,now)":{"date_range":{"field":"create_time","format":"yyyy-MM-dd","ranges":[{"from":"now-8d","to":"now-7d"},{"from":"now-7d","to":"now-6d"},{"from":"now-6d","to":"now"}]}}}}`,
			expectedTable: "mock_table",
			expectedErr:   "",
		},
		// Test cases with invalid SQL queries
		{
			sql:           "DELETE FROM mock_table",
			expectedDSL:   "",
			expectedTable: "",
			expectedErr:   "statement not supported",
		},
		{
			sql:           "SELECT * FROM mock_table",
			expectedDSL:   "",
			expectedTable: "",
			expectedErr:   "WHERE clause is missing",
		},
		{
			sql:           "SELECT * FROM mock_table, mock_table2 where id='123'",
			expectedDSL:   "",
			expectedTable: "",
			expectedErr:   "only one table supported",
		},
		{
			sql:           "SELECT * FROM mock_table WHERE name is null",
			expectedDSL:   "",
			expectedTable: "",
			expectedErr:   "is expression currently not supported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			dsl, table, err := Convert(tt.sql)

			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
				assert.Empty(t, dsl)
				assert.Empty(t, table)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDSL, dsl)
				assert.Equal(t, tt.expectedTable, table)
			}
		})
	}
}

func TestBuildComparisonExprRightStr(t *testing.T) {
	tests := []struct {
		name        string
		expr        sqlparser.Expr
		wantStr     string
		wantMissing bool
		wantErr     bool
	}{
		{
			"Integer value",
			&sqlparser.SQLVal{Type: sqlparser.IntVal, Val: []byte("123")},
			"123",
			false,
			false,
		},
		{
			"String value",
			&sqlparser.SQLVal{Type: sqlparser.StrVal, Val: []byte("'abc'")},
			`\'abc\`,
			false,
			false,
		},
		{
			"GroupConcat not supported",
			&sqlparser.GroupConcatExpr{}, // you might want to construct it fully
			"",
			false,
			true,
		},
		{
			"Function not supported",
			&sqlparser.FuncExpr{Name: sqlparser.NewColIdent("UnsupportedFunc")},
			"",
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rightStr, missingCheck, err := buildComparisonExprRightStr(tt.expr)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStr, rightStr)
				assert.Equal(t, tt.wantMissing, missingCheck)
			}
		})
	}
}

func TestHandleGroupByAgg(t *testing.T) {
	colName := sqlparser.NewColIdent("category")

	testCases := []struct {
		name      string
		groupBy   sqlparser.GroupBy
		innerMap  msi // msi as per your definition (map[string]interface{})
		wantAgg   msi
		wantError bool
	}{
		{
			"Single column groupBy",
			sqlparser.GroupBy{&sqlparser.ColName{Name: colName}},
			msi{}, // no inner aggregation map is given
			msi{
				"category": msi{
					"terms": msi{
						"field": "category",
						"size":  200,
					},
				},
			},
			false,
		},
		{
			"Multiple column groupBy",
			sqlparser.GroupBy{
				&sqlparser.ColName{Name: sqlparser.NewColIdent("category")},
				&sqlparser.ColName{Name: sqlparser.NewColIdent("sub_category")},
			},
			msi{}, // no inner aggregation map is given for simplicity, in real case it may be complex
			msi{
				"category": msi{
					"aggregations": msi{
						"sub_category": msi{
							"terms": msi{
								"field": "sub_category",
								"size":  0,
							},
						},
					}, "terms": msi{
						"field": "category",
						"size":  200,
					},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aggMap, err := handleGroupByAgg(tc.groupBy, tc.innerMap)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantAgg, aggMap)
			}
		})
	}
}
