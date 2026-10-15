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

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/KusionStack/karpor/pkg/infra/search/storage"
	"github.com/xwb1989/sqlparser"
)

const (
	allResourcesCondition = "1 = 1"
	maxSearchResults      = 1000000
)

type fieldExpression struct {
	sql      string
	args     []interface{}
	name     string
	isBool   bool
	nullable bool
}

func (s *Storage) Search(ctx context.Context, queryString, patternType string, pagination *storage.Pagination) (*storage.SearchResult, error) {
	switch patternType {
	case storage.SQLPatternType, storage.NLPatternType:
		plan, err := compileSearchQuery(queryString, pagination)
		if err != nil {
			return nil, err
		}
		return s.searchByPlan(ctx, plan)
	case storage.DSLPatternType:
		return nil, fmt.Errorf("sqlite storage does not support dsl search pattern")
	default:
		return nil, fmt.Errorf("invalid type %s", patternType)
	}
}

func (s *Storage) SearchByTerms(ctx context.Context, keysAndValues map[string]any, pagination *storage.Pagination) (*storage.SearchResult, error) {
	whereSQL, args, err := compileTermFilters(keysAndValues)
	if err != nil {
		return nil, err
	}
	limit, offset := paginationLimitOffset(pagination)
	return s.searchByPlan(ctx, &searchPlan{
		whereSQL:  whereSQL,
		whereArgs: args,
		orderSQL:  defaultResourceOrder(),
		limit:     limit,
		offset:    offset,
	})
}

func (s *Storage) AggregateByTerms(ctx context.Context, keys []string) (*storage.AggregateResults, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no fields provided for aggregation")
	}

	fields := make([]fieldExpression, 0, len(keys))
	selects := make([]string, 0, len(keys))
	where := make([]string, 0, len(keys))
	groupBy := make([]string, 0, len(keys))
	orderBy := make([]string, 0, len(keys))
	selectArgs := []interface{}{}
	whereArgs := []interface{}{}
	groupArgs := []interface{}{}
	for _, key := range keys {
		field, err := resolveField(key)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
		selects = append(selects, fmt.Sprintf("%s AS field_%d", field.sql, len(fields)-1))
		groupBy = append(groupBy, field.sql)
		orderBy = append(orderBy, fmt.Sprintf("field_%d", len(fields)-1))
		selectArgs = append(selectArgs, field.args...)
		groupArgs = append(groupArgs, field.args...)
		if field.nullable {
			where = append(where, field.sql+" IS NOT NULL")
			whereArgs = append(whereArgs, field.args...)
		}
	}

	query := fmt.Sprintf(`SELECT %s, COUNT(*) FROM resources`,
		strings.Join(selects, ", "),
	)
	args := append([]interface{}{}, selectArgs...)
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
		args = append(args, whereArgs...)
	}
	query += fmt.Sprintf(` GROUP BY %s ORDER BY COUNT(*) DESC, %s`,
		strings.Join(groupBy, ", "),
		strings.Join(orderBy, ", "),
	)
	args = append(args, groupArgs...)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buckets := []storage.Bucket{}
	for rows.Next() {
		values := make([]sql.NullString, len(keys))
		dest := make([]interface{}, 0, len(keys)+1)
		for i := range values {
			dest = append(dest, &values[i])
		}
		var count int
		dest = append(dest, &count)
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}

		bucket := storage.Bucket{Keys: make([]string, len(keys)), Count: count}
		for i, value := range values {
			if value.Valid {
				bucket.Keys[i] = value.String
			}
		}
		buckets = append(buckets, bucket)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &storage.AggregateResults{
		Buckets: buckets,
		Total:   len(buckets),
	}, nil
}

type searchPlan struct {
	whereSQL  string
	whereArgs []interface{}
	orderSQL  string
	orderArgs []interface{}
	limit     int
	offset    int
}

func (s *Storage) searchByPlan(ctx context.Context, plan *searchPlan) (*storage.SearchResult, error) {
	whereSQL := plan.whereSQL
	if whereSQL == "" {
		whereSQL = allResourcesCondition
	}

	countQuery := `SELECT COUNT(*) FROM resources WHERE ` + whereSQL
	var total int
	if err := s.db.QueryRowContext(ctx, countQuery, plan.whereArgs...).Scan(&total); err != nil {
		return nil, err
	}

	args := append([]interface{}{}, plan.whereArgs...)
	query := `SELECT uid, cluster, api_version, kind, namespace, name, content, sync_at, deleted
		FROM resources
		WHERE ` + whereSQL
	if plan.orderSQL != "" {
		query += ` ORDER BY ` + plan.orderSQL
		args = append(args, plan.orderArgs...)
	}
	if plan.limit > 0 {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, plan.limit, plan.offset)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resources := []*storage.Resource{}
	for rows.Next() {
		resource, _, err := scanResource(rows)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &storage.SearchResult{
		Total:     total,
		Resources: resources,
	}, nil
}

func compileSearchQuery(queryString string, pagination *storage.Pagination) (*searchPlan, error) {
	if strings.Contains(queryString, ";") {
		return nil, fmt.Errorf("multiple statements are not supported")
	}

	stmt, err := sqlparser.Parse(queryString)
	if err != nil {
		return nil, err
	}
	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, fmt.Errorf("statement not supported")
	}
	if err := validateSelect(sel); err != nil {
		return nil, err
	}

	var (
		whereSQL string
		args     []interface{}
	)
	if sel.Where != nil {
		whereSQL, args, err = compileExpr(sel.Where.Expr)
		if err != nil {
			return nil, err
		}
	}
	if !whereMentionsDeleted(sel.Where) {
		whereSQL, args = andCondition(whereSQL, args, "deleted = ?", []interface{}{0})
	}

	orderSQL, orderArgs, err := compileOrderBy(sel.OrderBy)
	if err != nil {
		return nil, err
	}
	if orderSQL == "" {
		orderSQL = defaultResourceOrder()
	}

	limit, offset := paginationLimitOffset(pagination)
	if pagination == nil && sel.Limit != nil {
		limit, offset, err = compileLimit(sel.Limit)
		if err != nil {
			return nil, err
		}
	}

	return &searchPlan{
		whereSQL:  whereSQL,
		whereArgs: args,
		orderSQL:  orderSQL,
		orderArgs: orderArgs,
		limit:     limit,
		offset:    offset,
	}, nil
}

func validateSelect(sel *sqlparser.Select) error {
	if len(sel.From) != 1 {
		return fmt.Errorf("only one table supported")
	}
	tableExpr, ok := sel.From[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return fmt.Errorf("joins and derived tables are not supported")
	}
	tableName, ok := tableExpr.Expr.(sqlparser.TableName)
	if !ok || !tableName.Qualifier.IsEmpty() || tableName.Name.String() != "resources" {
		return fmt.Errorf("only resources table is supported")
	}
	if !tableExpr.As.IsEmpty() {
		return fmt.Errorf("table aliases are not supported")
	}
	if len(sel.GroupBy) > 0 {
		return fmt.Errorf("GROUP BY is not supported in sqlite resource search")
	}
	if sel.Having != nil {
		return fmt.Errorf("HAVING is not supported in sqlite resource search")
	}
	for _, expr := range sel.SelectExprs {
		switch e := expr.(type) {
		case *sqlparser.StarExpr:
		case *sqlparser.AliasedExpr:
			if _, ok := e.Expr.(*sqlparser.ColName); !ok {
				return fmt.Errorf("only column selection is supported")
			}
		default:
			return fmt.Errorf("select expression is not supported")
		}
	}
	return nil
}

func compileExpr(expr sqlparser.Expr) (string, []interface{}, error) {
	switch e := expr.(type) {
	case *sqlparser.AndExpr:
		leftSQL, leftArgs, err := compileExpr(e.Left)
		if err != nil {
			return "", nil, err
		}
		rightSQL, rightArgs, err := compileExpr(e.Right)
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("(%s AND %s)", leftSQL, rightSQL), append(leftArgs, rightArgs...), nil
	case *sqlparser.OrExpr:
		leftSQL, leftArgs, err := compileExpr(e.Left)
		if err != nil {
			return "", nil, err
		}
		rightSQL, rightArgs, err := compileExpr(e.Right)
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("(%s OR %s)", leftSQL, rightSQL), append(leftArgs, rightArgs...), nil
	case *sqlparser.ParenExpr:
		innerSQL, args, err := compileExpr(e.Expr)
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("(%s)", innerSQL), args, nil
	case *sqlparser.ComparisonExpr:
		return compileComparison(e)
	case *sqlparser.RangeCond:
		return compileRange(e)
	case *sqlparser.FuncExpr:
		return compileFunction(e)
	default:
		return "", nil, fmt.Errorf("unsupported expression %T", expr)
	}
}

func compileComparison(expr *sqlparser.ComparisonExpr) (string, []interface{}, error) {
	colName, ok := expr.Left.(*sqlparser.ColName)
	if !ok {
		return "", nil, fmt.Errorf("invalid comparison expression, the left must be a column name")
	}
	field, err := resolveColName(colName)
	if err != nil {
		return "", nil, err
	}
	operator := strings.ToLower(expr.Operator)

	switch operator {
	case sqlparser.EqualStr, sqlparser.NotEqualStr, ">", "<", ">=", "<=":
		value, err := comparisonValue(expr.Right, field)
		if err != nil {
			return "", nil, err
		}
		switch operator {
		case sqlparser.EqualStr:
			return fmt.Sprintf("%s = ?", field.sql), append(field.args, value), nil
		case sqlparser.NotEqualStr:
			args := append(append([]interface{}{}, field.args...), field.args...)
			args = append(args, value)
			return fmt.Sprintf("(%s IS NULL OR %s != ?)", field.sql, field.sql), args, nil
		default:
			return fmt.Sprintf("%s %s ?", field.sql, operator), append(field.args, value), nil
		}
	case sqlparser.LikeStr:
		value, err := stringComparisonValue(expr.Right)
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("%s LIKE ? ESCAPE '\\'", field.sql), append(field.args, value), nil
	case sqlparser.NotLikeStr:
		value, err := stringComparisonValue(expr.Right)
		if err != nil {
			return "", nil, err
		}
		args := append(append([]interface{}{}, field.args...), field.args...)
		args = append(args, value)
		return fmt.Sprintf("(%s IS NULL OR %s NOT LIKE ? ESCAPE '\\')", field.sql, field.sql), args, nil
	case sqlparser.InStr, sqlparser.NotInStr:
		values, err := tupleValues(expr.Right, field)
		if err != nil {
			return "", nil, err
		}
		if len(values) == 0 {
			return "", nil, fmt.Errorf("IN requires at least one value")
		}
		placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
		if operator == sqlparser.InStr {
			return fmt.Sprintf("%s IN (%s)", field.sql, placeholders), append(field.args, values...), nil
		}
		args := append(append([]interface{}{}, field.args...), field.args...)
		args = append(args, values...)
		return fmt.Sprintf("(%s IS NULL OR %s NOT IN (%s))", field.sql, field.sql, placeholders), args, nil
	default:
		return "", nil, fmt.Errorf("unsupported comparison operator %s", expr.Operator)
	}
}

func compileRange(expr *sqlparser.RangeCond) (string, []interface{}, error) {
	colName, ok := expr.Left.(*sqlparser.ColName)
	if !ok {
		return "", nil, fmt.Errorf("range column name missing")
	}
	field, err := resolveColName(colName)
	if err != nil {
		return "", nil, err
	}
	from, err := comparisonValue(expr.From, field)
	if err != nil {
		return "", nil, err
	}
	to, err := comparisonValue(expr.To, field)
	if err != nil {
		return "", nil, err
	}
	switch strings.ToLower(expr.Operator) {
	case sqlparser.BetweenStr:
		return fmt.Sprintf("%s BETWEEN ? AND ?", field.sql), append(field.args, from, to), nil
	case sqlparser.NotBetweenStr:
		args := append(append([]interface{}{}, field.args...), field.args...)
		args = append(args, from, to)
		return fmt.Sprintf("(%s IS NULL OR %s NOT BETWEEN ? AND ?)", field.sql, field.sql), args, nil
	default:
		return "", nil, fmt.Errorf("unsupported range operator %s", expr.Operator)
	}
}

func compileFunction(expr *sqlparser.FuncExpr) (string, []interface{}, error) {
	switch expr.Name.Lowered() {
	case "contains":
		if len(expr.Exprs) != 2 {
			return "", nil, fmt.Errorf("contains expects 2 expressions")
		}
		fieldExpr, ok := expr.Exprs[0].(*sqlparser.AliasedExpr)
		if !ok {
			return "", nil, fmt.Errorf("contains field expression is not supported")
		}
		colName, ok := fieldExpr.Expr.(*sqlparser.ColName)
		if !ok {
			return "", nil, fmt.Errorf("contains field must be a column name")
		}
		field, err := resolveColName(colName)
		if err != nil {
			return "", nil, err
		}
		valueExpr, ok := expr.Exprs[1].(*sqlparser.AliasedExpr)
		if !ok {
			return "", nil, fmt.Errorf("contains value expression is not supported")
		}
		value, err := stringComparisonValue(valueExpr.Expr)
		if err != nil {
			return "", nil, err
		}
		return fmt.Sprintf("%s LIKE ? ESCAPE '\\'", field.sql), append(field.args, "%"+escapeLike(value)+"%"), nil
	default:
		return "", nil, fmt.Errorf("function in where not supported %s", expr.Name.Lowered())
	}
}

func compileTermFilters(keysAndValues map[string]any) (string, []interface{}, error) {
	if len(keysAndValues) == 0 {
		return allResourcesCondition, nil, nil
	}

	keys := make([]string, 0, len(keysAndValues))
	for key := range keysAndValues {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var whereSQL string
	var args []interface{}
	for _, key := range keys {
		conditionSQL, conditionArgs, err := compileTermFilter(key, keysAndValues[key])
		if err != nil {
			return "", nil, err
		}
		whereSQL, args = andCondition(whereSQL, args, conditionSQL, conditionArgs)
	}
	if whereSQL == "" {
		return allResourcesCondition, nil, nil
	}
	return whereSQL, args, nil
}

func compileTermFilter(key string, value any) (string, []interface{}, error) {
	if key == "labels" || key == "annotations" {
		return compileMapTermFilter(key, value)
	}
	field, err := resolveField(key)
	if err != nil {
		return "", nil, err
	}
	comparisonValue, err := normalizeValue(value, field)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf("%s = ?", field.sql), append(field.args, comparisonValue), nil
}

func compileMapTermFilter(prefix string, value any) (string, []interface{}, error) {
	values := map[string]string{}
	switch typed := value.(type) {
	case map[string]string:
		values = typed
	case map[string]interface{}:
		for key, val := range typed {
			strVal, ok := val.(string)
			if !ok {
				return "", nil, fmt.Errorf("unsupported %s value type %T", prefix, val)
			}
			values[key] = strVal
		}
	default:
		return "", nil, fmt.Errorf("unsupported %s terms type %T", prefix, value)
	}
	if len(values) == 0 {
		return allResourcesCondition, nil, nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var whereSQL string
	var args []interface{}
	for _, key := range keys {
		field, err := resolveField(prefix + "." + key)
		if err != nil {
			return "", nil, err
		}
		whereSQL, args = andCondition(whereSQL, args, fmt.Sprintf("%s = ?", field.sql), append(field.args, values[key]))
	}
	return whereSQL, args, nil
}

func resolveColName(colName *sqlparser.ColName) (fieldExpression, error) {
	if !colName.Qualifier.IsEmpty() {
		return fieldExpression{}, fmt.Errorf("qualified columns are not supported")
	}
	return resolveField(colName.Name.String())
}

func resolveField(name string) (fieldExpression, error) {
	switch name {
	case "uid":
		return fieldExpression{sql: "uid", name: name}, nil
	case "cluster":
		return fieldExpression{sql: "cluster", name: name}, nil
	case "apiVersion":
		return fieldExpression{sql: "api_version", name: name}, nil
	case "kind":
		return fieldExpression{sql: "kind", name: name}, nil
	case "namespace":
		return fieldExpression{sql: "namespace", name: name}, nil
	case "name":
		return fieldExpression{sql: "name", name: name}, nil
	case "creationTimestamp":
		return fieldExpression{sql: "creation_timestamp", name: name}, nil
	case "deletionTimestamp":
		return fieldExpression{sql: "deletion_timestamp", name: name}, nil
	case "ownerReferences":
		return fieldExpression{sql: "owner_references", name: name}, nil
	case "resourceVersion":
		return fieldExpression{sql: "resource_version", name: name}, nil
	case "content":
		return fieldExpression{sql: "content", name: name}, nil
	case "syncAt":
		return fieldExpression{sql: "sync_at", name: name}, nil
	case "deleted":
		return fieldExpression{sql: "deleted", name: name, isBool: true}, nil
	default:
		switch {
		case strings.HasPrefix(name, "labels."):
			path, err := jsonPath(strings.TrimPrefix(name, "labels."))
			if err != nil {
				return fieldExpression{}, err
			}
			return fieldExpression{sql: "json_extract(labels, ?)", args: []interface{}{path}, name: name, nullable: true}, nil
		case strings.HasPrefix(name, "annotations."):
			path, err := jsonPath(strings.TrimPrefix(name, "annotations."))
			if err != nil {
				return fieldExpression{}, err
			}
			return fieldExpression{sql: "json_extract(annotations, ?)", args: []interface{}{path}, name: name, nullable: true}, nil
		default:
			return fieldExpression{}, fmt.Errorf("unsupported field %s", name)
		}
	}
}

func comparisonValue(expr sqlparser.Expr, field fieldExpression) (interface{}, error) {
	switch value := expr.(type) {
	case *sqlparser.SQLVal:
		return normalizeValue(string(value.Val), field)
	case sqlparser.BoolVal:
		return normalizeValue(bool(value), field)
	default:
		return nil, fmt.Errorf("only literal comparison values are supported")
	}
}

func stringComparisonValue(expr sqlparser.Expr) (string, error) {
	sqlVal, ok := expr.(*sqlparser.SQLVal)
	if !ok || sqlVal.Type != sqlparser.StrVal {
		return "", fmt.Errorf("only string literal values are supported")
	}
	return string(sqlVal.Val), nil
}

func tupleValues(expr sqlparser.Expr, field fieldExpression) ([]interface{}, error) {
	tuple, ok := expr.(sqlparser.ValTuple)
	if !ok {
		return nil, fmt.Errorf("IN requires a value tuple")
	}
	values := make([]interface{}, 0, len(tuple))
	for _, tupleExpr := range tuple {
		value, err := comparisonValue(tupleExpr, field)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func normalizeValue(value any, field fieldExpression) (interface{}, error) {
	if !field.isBool {
		return value, nil
	}
	switch typed := value.(type) {
	case bool:
		if typed {
			return 1, nil
		}
		return 0, nil
	case string:
		switch strings.ToLower(typed) {
		case "true", "1":
			return 1, nil
		case "false", "0":
			return 0, nil
		default:
			return nil, fmt.Errorf("invalid boolean value %q for field %s", typed, field.name)
		}
	case int:
		if typed != 0 {
			return 1, nil
		}
		return 0, nil
	default:
		return nil, fmt.Errorf("invalid boolean value type %T for field %s", value, field.name)
	}
}

func compileOrderBy(orderBy sqlparser.OrderBy) (string, []interface{}, error) {
	parts := make([]string, 0, len(orderBy))
	args := []interface{}{}
	for _, expr := range orderBy {
		colName, ok := expr.Expr.(*sqlparser.ColName)
		if !ok {
			return "", nil, fmt.Errorf("only column ORDER BY is supported")
		}
		field, err := resolveColName(colName)
		if err != nil {
			return "", nil, err
		}
		direction := strings.ToLower(expr.Direction)
		if direction == "" {
			direction = sqlparser.AscScr
		}
		if direction != sqlparser.AscScr && direction != sqlparser.DescScr {
			return "", nil, fmt.Errorf("unsupported ORDER BY direction %s", expr.Direction)
		}
		parts = append(parts, fmt.Sprintf("%s %s", field.sql, direction))
		args = append(args, field.args...)
	}
	return strings.Join(parts, ", "), args, nil
}

func compileLimit(limit *sqlparser.Limit) (int, int, error) {
	rowCount, err := intLiteral(limit.Rowcount)
	if err != nil {
		return 0, 0, err
	}
	offset := 0
	if limit.Offset != nil {
		offset, err = intLiteral(limit.Offset)
		if err != nil {
			return 0, 0, err
		}
	}
	return rowCount, offset, nil
}

func intLiteral(expr sqlparser.Expr) (int, error) {
	sqlVal, ok := expr.(*sqlparser.SQLVal)
	if !ok || sqlVal.Type != sqlparser.IntVal {
		return 0, fmt.Errorf("limit values must be integer literals")
	}
	value, err := strconv.Atoi(string(sqlVal.Val))
	if err != nil {
		return 0, err
	}
	if value < 0 {
		return 0, fmt.Errorf("limit values cannot be negative")
	}
	return value, nil
}

func paginationLimitOffset(pagination *storage.Pagination) (int, int) {
	if pagination == nil || pagination.PageSize <= 0 {
		return maxSearchResults, 0
	}
	page := pagination.Page
	if page <= 1 {
		page = 1
	}
	return pagination.PageSize, (page - 1) * pagination.PageSize
}

func whereMentionsDeleted(where *sqlparser.Where) bool {
	if where == nil {
		return false
	}
	found := false
	where.Expr.WalkSubtree(func(node sqlparser.SQLNode) (bool, error) {
		if colName, ok := node.(*sqlparser.ColName); ok && colName.Name.String() == "deleted" {
			found = true
			return false, nil
		}
		return true, nil
	})
	return found
}

func andCondition(existing string, existingArgs []interface{}, condition string, conditionArgs []interface{}) (string, []interface{}) {
	if existing == "" {
		return condition, append(existingArgs, conditionArgs...)
	}
	return fmt.Sprintf("(%s AND %s)", existing, condition), append(existingArgs, conditionArgs...)
}

func jsonPath(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("json field key cannot be empty")
	}
	key = strings.ReplaceAll(key, `\`, `\\`)
	key = strings.ReplaceAll(key, `"`, `\"`)
	return `$."` + key + `"`, nil
}

func escapeLike(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `%`, `\%`)
	value = strings.ReplaceAll(value, `_`, `\_`)
	return value
}

func defaultResourceOrder() string {
	return "cluster ASC, api_version ASC, kind ASC, namespace ASC, name ASC, uid ASC"
}
