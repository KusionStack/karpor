/**
 * @Author:      Adam wu
 * @Description:
 * @File:        convert.go
 * @Version:     1.0.0
 * @Date:        2025/3/21
 */

package meilisearch

import (
	"fmt"
	"github.com/KusionStack/karpor/pkg/infra/persistence/meilisearch"
	"github.com/xwb1989/sqlparser"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
	"time"
)

var DefaultFilter = []string{"deleted = false"}

// Convert will convert a SQL query string to a MeiliSearch SearchRequest.
// It will also return the table name extracted from the SQL query.
// But there is some limitations:
// 1. MeiliSearch do not support group by.
// 2. MeiliSearch do not support Left Like, all Left Like will be replaced by CONTAINS.
// 3. MeiliSearch do not support string compare, so SyncAt should translate to unix timestamp .
func Convert(sql string) (*meilisearch.SearchRequest, string, error) {
	return ConvertWithDefaultFilter(sql, DefaultFilter)
}

func ConvertWithDefaultFilter(sql string, defaultFilter []string) (*meilisearch.SearchRequest, string, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, "", err
	}

	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, "", fmt.Errorf("only SELECT statements are supported")
	}
	if sel.Where == nil {
		return nil, "", fmt.Errorf("WHERE clause is required")
	}
	if len(sel.From) != 1 {
		return nil, "", fmt.Errorf("only single table queries are supported")
	}

	// 应用默认过滤器
	sel = applyDefaultFilter(sel, defaultFilter)

	// 获取表名
	tableName := strings.ReplaceAll(sqlparser.String(sel.From), "`", "")

	// 构建搜索请求
	req := &meilisearch.SearchRequest{
		Limit:  1000, //by default
		Offset: 0,
	}

	// 处理WHERE条件
	filter, err := buildFilter(sel.Where)
	if err != nil {
		return nil, "", err
	}
	req.Filter = filter

	// 处理LIMIT/OFFSET
	if sel.Limit != nil {
		if offset, ok := sel.Limit.Offset.(*sqlparser.SQLVal); ok {
			req.Offset = parseInt(string(offset.Val))
		}
		if limit, ok := sel.Limit.Rowcount.(*sqlparser.SQLVal); ok {
			req.Limit = parseInt(string(limit.Val))
		}
	}

	// 处理ORDER BY
	if len(sel.OrderBy) > 0 {
		sort := make([]string, 0, len(sel.OrderBy))
		for _, order := range sel.OrderBy {
			field := strings.Trim(sqlparser.String(order.Expr), "`")
			direction := strings.ToLower(order.Direction)
			sort = append(sort, fmt.Sprintf("%s:%s", field, direction))
		}
		req.Sort = sort
	}

	return req, tableName, nil
}

func buildFilter(where *sqlparser.Where) (interface{}, error) {
	return buildFilterRecursive(where.Expr)
}

func buildFilterRecursive(expr sqlparser.Expr) (interface{}, error) {
	switch e := expr.(type) {
	case *sqlparser.AndExpr:
		left, err := buildFilterRecursive(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := buildFilterRecursive(e.Right)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("(%s) AND (%s)", left, right), nil

	case *sqlparser.OrExpr:
		left, err := buildFilterRecursive(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := buildFilterRecursive(e.Right)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("(%s) OR (%s)", left, right), nil

	case *sqlparser.ComparisonExpr:
		return buildComparisonFilter(e)

	case *sqlparser.ParenExpr:
		return buildFilterRecursive(e.Expr)

	case *sqlparser.RangeCond:
		return buildRangeFilter(e)

	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func buildComparisonFilter(expr *sqlparser.ComparisonExpr) (string, error) {
	field := strings.Trim(sqlparser.String(expr.Left), "`")
	op, value := extractOperatorAndValue(field, expr.Operator, expr.Right)
	return fmt.Sprintf("%s %s %s", field, op, value), nil
}

func buildRangeFilter(expr *sqlparser.RangeCond) (string, error) {
	col, ok := expr.Left.(*sqlparser.ColName)
	if !ok {
		return "", fmt.Errorf("invalid range condition")
	}

	field := strings.Trim(sqlparser.String(col), "`")
	from := extractValue(expr.From)
	to := extractValue(expr.To)
	if expr.Operator == sqlparser.BetweenStr {
		return fmt.Sprintf("(%s >= %s) AND (%s <= %s)", field, from, field, to), nil
	}
	if expr.Operator == sqlparser.NotBetweenStr {
		return fmt.Sprintf("(%s < %s) OR (%s > %s)", field, from, field, to), nil
	}
	return "", fmt.Errorf("invalid range condition")
}

func extractLikeOperatorByValue(op, val string) string {
	val = strings.Trim(val, "'")
	leftLike := strings.HasPrefix(val, "%")
	rightLike := strings.HasSuffix(val, "%")
	if op == sqlparser.LikeStr {
		if leftLike {
			return "CONTAINS"
		}
		if rightLike {
			return "STARTS WITH"
		}
	}
	if op == sqlparser.NotLikeStr {
		if leftLike {
			return "NOT CONTAINS"
		}
		if rightLike {
			return "NOT STARTS WITH"
		}
	}
	return "="
}

// trimLikeValue trim string like '%abc%' to 'abc'
func trimLikeValue(val string) string {

	bs := []byte(val)
	newBytes := make([]byte, 0, len(bs))
	for i := 0; i < len(bs); i++ {
		if (i <= 1 || i >= len(bs)-2) && bs[i] == '%' {
			continue
		}
		newBytes = append(newBytes, bs[i])
	}
	return string(newBytes)
}

func extractOperatorAndValue(field, op string, expr sqlparser.Expr) (string, string) {
	op = strings.ToLower(op)
	val := extractValue(expr)
	switch op {
	//case sqlparser.EqualStr, sqlparser.GreaterEqualStr, sqlparser.GreaterThanStr, sqlparser.LessThanStr, sqlparser.LessEqualStr, sqlparser.NotEqualStr, sqlparser.InStr, sqlparser.NotInStr:
	//	return strings.ToUpper(op), val
	case sqlparser.LikeStr, sqlparser.NotLikeStr:
		return extractLikeOperatorByValue(op, val), trimLikeValue(val)
	case sqlparser.GreaterEqualStr, sqlparser.GreaterThanStr, sqlparser.LessThanStr, sqlparser.LessEqualStr:
		// MeiliSearch does not support string comparison, so we need to convert it to unix timestamp
		if field == resourceGroupRuleKeyCreatedAt || field == resourceGroupRuleKeyUpdatedAt || field == resourceGroupRuleKeyDeletedAt || field == resourceKeySyncAt {
			val = strings.Trim(val, "'")
			parse, _ := time.Parse(time.RFC3339, val)
			val = fmt.Sprintf("%d", parse.Unix())
		}
		return strings.ToUpper(op), val
	default:
		return strings.ToUpper(op), val
	}
}

func extractValue(expr sqlparser.Expr) string {
	switch v := expr.(type) {
	case *sqlparser.SQLVal:
		if v.Type == sqlparser.StrVal {
			return fmt.Sprintf("'%s'", string(v.Val))
		}
		return string(v.Val)
	case sqlparser.ValTuple:
		buf := sqlparser.NewTrackedBuffer(nil)
		buf.Myprintf("[%v]", sqlparser.Exprs(v))
		return buf.String()
	default:
		return sqlparser.String(expr)
	}
}

func applyDefaultFilter(sel *sqlparser.Select, filter []string) *sqlparser.Select {
	if len(filter) == 0 {
		return sel
	}

	existingFields := getFilterFields(sel.Where)
	defaultFields := sets.New[string]()
	for _, f := range filter {
		parts := strings.Split(f, " ")
		if len(parts) > 0 {
			defaultFields.Insert(parts[0])
		}
	}

	if existingFields.Intersection(defaultFields).Len() == 0 {
		for _, f := range filter {
			sel.AddWhere(&sqlparser.ComparisonExpr{
				Operator: "=",
				Left:     &sqlparser.ColName{Name: sqlparser.NewColIdent(strings.Split(f, " ")[0])},
				Right:    sqlparser.NewStrVal([]byte(strings.Join(strings.Split(f, " ")[2:], " "))),
			})
		}
	}

	return sel
}

func getFilterFields(where *sqlparser.Where) sets.Set[string] {
	fields := sets.New[string]()
	if where == nil {
		return fields
	}

	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (bool, error) {
		switch n := node.(type) {
		case *sqlparser.ColName:
			fields.Insert(n.Name.String())
		}
		return true, nil
	}, where)

	return fields
}

func parseInt(s string) int64 {
	var n int64
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}
