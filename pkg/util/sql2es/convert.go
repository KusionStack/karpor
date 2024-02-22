package sql2es

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func Convert(sql string) (dsl string, table string, err error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return "", "", err
	}

	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return "", "", fmt.Errorf("statement not supported")
	}
	if sel.Where == nil {
		return "", "", fmt.Errorf("WHERE clause is missing")
	}
	if len(sel.From) != 1 {
		return "", "", fmt.Errorf("only one table supported")
	}

	return handleSelect(sel)
}

func handleSelect(sel *sqlparser.Select) (dsl string, esType string, err error) {
	var rootParent sqlparser.Expr

	queryMapStr, err := handleSelectWhere(&sel.Where.Expr, true, &rootParent)
	if err != nil {
		return "", "", err
	} else if queryMapStr == "" {
		return "", "", fmt.Errorf("failed to generate query")
	}

	esType = sqlparser.String(sel.From)
	esType = strings.Replace(esType, "`", "", -1)

	queryFrom, querySize := "0", "1"

	aggFlag := false
	// if the request is to aggregation
	// then set aggFlag to true, and querySize to 0
	// to not return any query result

	var aggStr string
	if len(sel.GroupBy) > 0 || checkNeedAgg(sel.SelectExprs) {
		aggFlag = true
		querySize = "0"
		aggStr, err = buildAggs(sel)
		if err != nil {
			// aggStr = ""
			return "", "", err
		}
	}

	// Handle limit
	if sel.Limit != nil {
		if sel.Limit.Offset != nil {
			queryFrom = sqlparser.String(sel.Limit.Offset)
		}
		querySize = sqlparser.String(sel.Limit.Rowcount)
	}

	// Handle order by
	// when executating aggregations, order by is useless
	var orderByArr []string
	if aggFlag == false {
		for _, orderByExpr := range sel.OrderBy {
			orderByStr := fmt.Sprintf(`{"%v": "%v"}`, strings.Replace(sqlparser.String(orderByExpr.Expr), "`", "", -1), orderByExpr.Direction)
			orderByArr = append(orderByArr, orderByStr)
		}
	}

	resultMap := map[string]interface{}{
		"query": queryMapStr,
		"from":  queryFrom,
		"size":  querySize,
	}

	if len(aggStr) > 0 {
		resultMap["aggregations"] = aggStr
	}

	if len(orderByArr) > 0 {
		resultMap["sort"] = fmt.Sprintf("[%v]", strings.Join(orderByArr, ","))
	}

	// keep the travesal in order, avoid unpredicted json
	keySlice := []string{"query", "from", "size", "sort", "aggregations"}
	var resultArr []string
	for _, mapKey := range keySlice {
		if val, ok := resultMap[mapKey]; ok {
			resultArr = append(resultArr, fmt.Sprintf(`"%v" : %v`, mapKey, val))
		}
	}

	dsl = "{" + strings.Join(resultArr, ",") + "}"
	return dsl, esType, nil
}

// if the where is empty, need to check whether to agg or not
func checkNeedAgg(sqlSelect sqlparser.SelectExprs) bool {
	for _, v := range sqlSelect {
		expr, ok := v.(*sqlparser.AliasedExpr)
		if !ok {
			// no need to handle, star expression * just skip is ok
			continue
		}

		// TODO more precise
		if _, ok := expr.Expr.(*sqlparser.FuncExpr); ok {
			return true
		}
	}
	return false
}

func buildNestedFuncStrValue(nestedFunc *sqlparser.FuncExpr) (string, error) {
	return "", fmt.Errorf("elasticsql: unsupported function" + nestedFunc.Name.String())
}

func handleSelectWhereAndExpr(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, error) {
	andExpr := (*expr).(*sqlparser.AndExpr)
	leftExpr := andExpr.Left
	rightExpr := andExpr.Right
	leftStr, err := handleSelectWhere(&leftExpr, false, expr)
	if err != nil {
		return "", err
	}
	rightStr, err := handleSelectWhere(&rightExpr, false, expr)
	if err != nil {
		return "", err
	}

	// not toplevel
	// if the parent node is also and, then the result can be merged

	var resultStr string
	if leftStr == "" || rightStr == "" {
		resultStr = leftStr + rightStr
	} else {
		resultStr = leftStr + `,` + rightStr
	}

	if _, ok := (*parent).(*sqlparser.AndExpr); ok {
		return resultStr, nil
	}
	return fmt.Sprintf(`{"bool" : {"must" : [%v]}}`, resultStr), nil
}

func handleSelectWhereOrExpr(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, error) {
	orExpr := (*expr).(*sqlparser.OrExpr)
	leftExpr := orExpr.Left
	rightExpr := orExpr.Right

	leftStr, err := handleSelectWhere(&leftExpr, false, expr)
	if err != nil {
		return "", err
	}

	rightStr, err := handleSelectWhere(&rightExpr, false, expr)
	if err != nil {
		return "", err
	}

	var resultStr string
	if leftStr == "" || rightStr == "" {
		resultStr = leftStr + rightStr
	} else {
		resultStr = leftStr + `,` + rightStr
	}

	// not toplevel
	// if the parent node is also or node, then merge the query param
	if _, ok := (*parent).(*sqlparser.OrExpr); ok {
		return resultStr, nil
	}

	return fmt.Sprintf(`{"bool" : {"should" : [%v]}}`, resultStr), nil
}

func buildComparisonExprRightStr(expr sqlparser.Expr) (string, bool, error) {
	var rightStr string
	var err error
	missingCheck := false
	switch expr.(type) {
	case *sqlparser.SQLVal:
		rightStr = sqlparser.String(expr)
		rightStr = strings.Trim(rightStr, `'`)
	case *sqlparser.GroupConcatExpr:
		return "", missingCheck, fmt.Errorf("elasticsql: group_concat not supported")
	case *sqlparser.FuncExpr:
		// parse nested
		funcExpr := expr.(*sqlparser.FuncExpr)
		rightStr, err = buildNestedFuncStrValue(funcExpr)
		if err != nil {
			return "", missingCheck, err
		}
	case *sqlparser.ColName:
		if strings.ToLower(sqlparser.String(expr)) == "missing" {
			missingCheck = true
			return "", missingCheck, nil
		}

		return "", missingCheck, fmt.Errorf("elasticsql: column name on the right side of compare operator is not supported")
	case sqlparser.ValTuple:
		rightStr = sqlparser.String(expr)
	default:
		// cannot reach here
	}
	return rightStr, missingCheck, err
}

func unescapeSql(sql, escapeStr string) string {
	resSql := ""
	strSegments := strings.Split(sql, escapeStr)
	for _, segment := range strSegments {
		if segment == "" {
			resSql += escapeStr // continuous escapeStr, only remove the first one
		} else {
			resSql += segment
		}
	}
	return resSql
}

func handleSelectWhereComparisonExpr(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, error) {
	comparisonExpr := (*expr).(*sqlparser.ComparisonExpr)
	colName, ok := comparisonExpr.Left.(*sqlparser.ColName)

	if !ok {
		return "", fmt.Errorf("elasticsql: invalid comparison expression, the left must be a column name")
	}

	colNameStr := sqlparser.String(colName)
	colNameStr = strings.Replace(colNameStr, "`", "", -1)
	rightStr, missingCheck, err := buildComparisonExprRightStr(comparisonExpr.Right)
	if err != nil {
		return "", err
	}

	// unescape rightStr
	escapeStr := sqlparser.String(comparisonExpr.Escape)
	escapeStr = strings.Trim(escapeStr, "'") // remove quote both sides
	rightStr = unescapeSql(rightStr, escapeStr)
	resultStr := ""

	// allow eq empty string
	if rightStr == `<nil>` {
		rightStr = ""
	}

	switch comparisonExpr.Operator {
	case ">=":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"from" : "%v"}}}`, colNameStr, rightStr)
	case "<=":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"to" : "%v"}}}`, colNameStr, rightStr)
	case "=":
		// field is missing
		if missingCheck { // missing was deprecated in 2.2, use exists instead.
			resultStr = fmt.Sprintf(`{"bool" : {"must_not" : [{"exists":{"field":"%v"}}]}}`, colNameStr)
		} else {
			resultStr = fmt.Sprintf(`{"match_phrase" : {"%v" : {"query" : "%v"}}}`, colNameStr, rightStr)
		}
	case ">":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"gt" : "%v"}}}`, colNameStr, rightStr)
	case "<":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"lt" : "%v"}}}`, colNameStr, rightStr)
	case "!=":
		if missingCheck { // missing was deprecated in 2.2, use exists instead.
			resultStr = fmt.Sprintf(`{"bool" : {"must" : [{"exists":{"field":"%v"}}]}}`, colNameStr)
		} else {
			resultStr = fmt.Sprintf(`{"bool" : {"must_not" : [{"match_phrase" : {"%v" : {"query" : "%v"}}}]}}`, colNameStr, rightStr)
		}
	case "in":
		// the default valTuple is ('1', '2', '3') like
		// so need to drop the () and replace ' to "
		rightStr = strings.Replace(rightStr, `'`, `"`, -1)
		rightStr = strings.Trim(rightStr, "(")
		rightStr = strings.Trim(rightStr, ")")
		resultStr = fmt.Sprintf(`{"terms" : {"%v" : [%v]}}`, colNameStr, rightStr)
	case "like":
		// rightStr = strings.Replace(rightStr, `%`, ``, -1)
		// resultStr = fmt.Sprintf(`{"match_phrase" : {"%v" : {"query" : "%v"}}}`, colNameStr, rightStr)
		rightStr = strings.Replace(rightStr, `%`, `*`, -1)
		rightStr = strings.Replace(rightStr, `_`, `?`, -1)
		resultStr = fmt.Sprintf(`{"wildcard" : {"%v" : "%v"}}`, colNameStr, rightStr)
	case "not like":
		// rightStr = strings.Replace(rightStr, `%`, ``, -1)
		// resultStr = fmt.Sprintf(`{"bool" : {"must_not" : {"match_phrase" : {"%v" : {"query" : "%v"}}}}}`, colNameStr, rightStr)
		rightStr = strings.Replace(rightStr, `%`, `*`, -1)
		rightStr = strings.Replace(rightStr, `_`, `?`, -1)
		resultStr = fmt.Sprintf(`{"bool" : {"must_not" : {"wildcard" : {"%v" : "%v"}}}}`, colNameStr, rightStr)
	case "not in":
		// the default valTuple is ('1', '2', '3') like
		// so need to drop the () and replace ' to "
		rightStr = strings.Replace(rightStr, `'`, `"`, -1)
		rightStr = strings.Trim(rightStr, "(")
		rightStr = strings.Trim(rightStr, ")")
		resultStr = fmt.Sprintf(`{"bool" : {"must_not" : {"terms" : {"%v" : [%v]}}}}`, colNameStr, rightStr)
	}

	// the root node need to have bool and must
	if topLevel {
		resultStr = fmt.Sprintf(`{"bool" : {"must" : [%v]}}`, resultStr)
	}

	return resultStr, nil
}

func handleSelectWhere(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, error) {
	if expr == nil {
		return "", fmt.Errorf("elasticsql: error expression cannot be nil here")
	}

	switch e := (*expr).(type) {
	case *sqlparser.AndExpr:
		return handleSelectWhereAndExpr(expr, topLevel, parent)

	case *sqlparser.OrExpr:
		return handleSelectWhereOrExpr(expr, topLevel, parent)
	case *sqlparser.ComparisonExpr:
		return handleSelectWhereComparisonExpr(expr, topLevel, parent)

	case *sqlparser.IsExpr:
		return "", fmt.Errorf("elasticsql: is expression currently not supported")
	case *sqlparser.RangeCond:
		// between a and b
		// the meaning is equal to range query
		rangeCond := (*expr).(*sqlparser.RangeCond)
		colName, ok := rangeCond.Left.(*sqlparser.ColName)

		if !ok {
			return "", fmt.Errorf("elasticsql: range column name missing")
		}

		colNameStr := sqlparser.String(colName)
		colNameStr = strings.Replace(colNameStr, "`", "", -1)
		fromStr := strings.Trim(sqlparser.String(rangeCond.From), `'`)
		toStr := strings.Trim(sqlparser.String(rangeCond.To), `'`)

		resultStr := fmt.Sprintf(`{"range" : {"%v" : {"from" : "%v", "to" : "%v"}}}`, colNameStr, fromStr, toStr)
		if topLevel {
			resultStr = fmt.Sprintf(`{"bool" : {"must" : [%v]}}`, resultStr)
		}

		return resultStr, nil

	case *sqlparser.ParenExpr:
		parentBoolExpr := (*expr).(*sqlparser.ParenExpr)
		boolExpr := parentBoolExpr.Expr

		// if paren is the top level, bool must is needed
		isThisTopLevel := false
		if topLevel {
			isThisTopLevel = true
		}
		return handleSelectWhere(&boolExpr, isThisTopLevel, parent)
	case *sqlparser.NotExpr:
		return "", fmt.Errorf("elasticsql: not expression currently not supported")
	case *sqlparser.FuncExpr:
		switch e.Name.Lowered() {
		case "multi_match":
			params := e.Exprs
			if len(params) > 3 || len(params) < 2 {
				return "", fmt.Errorf("elasticsql: the multi_match must have 2 or 3 params, (query, fields and type) or (query, fields)")
			}

			var typ, query, fields string
			for i := 0; i < len(params); i++ {
				elem := strings.Replace(sqlparser.String(params[i]), "`", "", -1) // a = b
				kv := strings.Split(elem, "=")
				if len(kv) != 2 {
					return "", fmt.Errorf("elasticsql: the param should be query = xxx, field = yyy, type = zzz")
				}
				k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
				switch k {
				case "type":
					typ = strings.Replace(v, "'", "", -1)
				case "query":
					query = strings.Replace(v, "`", "", -1)
					query = strings.Replace(query, "'", "", -1)
				case "fields":
					fieldList := strings.Split(strings.TrimRight(strings.TrimLeft(v, "("), ")"), ",")
					for idx, field := range fieldList {
						fieldList[idx] = fmt.Sprintf(`"%v"`, strings.TrimSpace(field))
					}
					fields = strings.Join(fieldList, ",")
				default:
					return "", fmt.Errorf("elaticsql: unknow param for multi_match")
				}
			}
			if typ == "" {
				return fmt.Sprintf(`{"multi_match" : {"query" : "%v", "fields" : [%v]}}`, query, fields), nil
			}
			return fmt.Sprintf(`{"multi_match" : {"query" : "%v", "type" : "%v", "fields" : [%v]}}`, query, typ, fields), nil
		default:
			return "", fmt.Errorf("elaticsql: function in where not supported" + e.Name.Lowered())
		}
	}

	return "", fmt.Errorf("elaticsql: logically cannot reached here")
}

// msi stands for map[string]interface{}
type msi map[string]interface{}

func handleFuncInSelectAgg(funcExprArr []*sqlparser.FuncExpr) msi {
	innerAggMap := make(msi)
	for _, v := range funcExprArr {
		// func expressions will use the same parent bucket

		aggName := strings.ToUpper(v.Name.String()) + `(` + sqlparser.String(v.Exprs) + `)`
		switch v.Name.Lowered() {
		case "count":
			// count need to distinguish * and normal field name
			if sqlparser.String(v.Exprs) == "*" {
				innerAggMap[aggName] = msi{
					"value_count": msi{
						"field": "_index",
					},
				}
			} else {
				// support count(distinct field)
				if v.Distinct {
					innerAggMap[aggName] = msi{
						"cardinality": msi{
							"field": sqlparser.String(v.Exprs),
						},
					}
				} else {
					innerAggMap[aggName] = msi{
						"value_count": msi{
							"field": sqlparser.String(v.Exprs),
						},
					}
				}
			}
		default:
			// support min/avg/max/stats
			// extended_stats/percentiles
			innerAggMap[aggName] = msi{
				v.Name.String(): msi{
					"field": sqlparser.String(v.Exprs),
				},
			}
		}

	}

	return innerAggMap
}

func handleGroupByColName(colName *sqlparser.ColName, index int, child msi) msi {
	innerMap := make(msi)
	if index == 0 {
		innerMap["terms"] = msi{
			"field": colName.Name.String(),
			"size":  200, // this size may need to change ?
		}
	} else {
		innerMap["terms"] = msi{
			"field": colName.Name.String(),
			"size":  0,
		}
	}

	if len(child) > 0 {
		innerMap["aggregations"] = child
	}
	return msi{colName.Name.String(): innerMap}
}

func handleGroupByFuncExprDateHisto(funcExpr *sqlparser.FuncExpr) (msi, error) {
	innerMap := make(msi)
	var (
		// default
		field    = ""
		interval = "1h"
		format   = "yyyy-MM-dd HH:mm:ss"
	)

	// get field/interval and format
	for _, expr := range funcExpr.Exprs {
		// the expression in date_histogram must be like a = b format
		switch item := expr.(type) {
		case *sqlparser.AliasedExpr:
			// nonStarExpr := expr.(*sqlparser.NonStarExpr)
			comparisonExpr, ok := item.Expr.(*sqlparser.ComparisonExpr)

			if !ok {
				return nil, fmt.Errorf("elasticsql: unsupported expression in date_histogram")
			}
			left, ok := comparisonExpr.Left.(*sqlparser.ColName)
			if !ok {
				return nil, fmt.Errorf("elaticsql: param error in date_histogram")
			}
			rightStr := sqlparser.String(comparisonExpr.Right)
			rightStr = strings.Replace(rightStr, `'`, ``, -1)
			if left.Name.Lowered() == "field" {
				field = rightStr
			}
			if left.Name.Lowered() == "value" || left.Name.Lowered() == "interval" {
				interval = rightStr
			}
			if left.Name.Lowered() == "format" {
				format = rightStr
			}

			innerMap["date_histogram"] = msi{
				"field":    field,
				"interval": interval,
				"format":   format,
			}
		default:
			return nil, fmt.Errorf("elasticsql: unsupported expression in date_histogram")
		}
	}
	return innerMap, nil
}

func handleGroupByFuncExprRange(funcExpr *sqlparser.FuncExpr) (msi, error) {
	if len(funcExpr.Exprs) < 3 {
		return nil, fmt.Errorf("elasticsql: length of function range params must be > 3")
	}

	innerMap := make(msi)
	rangeMapList := make([]msi, len(funcExpr.Exprs)-2)

	for i := 1; i < len(funcExpr.Exprs)-1; i++ {
		valFrom := sqlparser.String(funcExpr.Exprs[i])
		valTo := sqlparser.String(funcExpr.Exprs[i+1])
		rangeMapList[i-1] = msi{
			"from": valFrom,
			"to":   valTo,
		}
	}
	innerMap[funcExpr.Name.String()] = msi{
		"field":  sqlparser.String(funcExpr.Exprs[0]),
		"ranges": rangeMapList,
	}

	return innerMap, nil
}

func handleGroupByFuncExprDateRange(funcExpr *sqlparser.FuncExpr) (msi, error) {
	var innerMap msi
	var (
		field        string
		format       = "yyyy-MM-dd HH:mm:ss"
		rangeList    = []string{}
		rangeMapList = []msi{}
	)

	for _, expr := range funcExpr.Exprs {
		nonStarExpr, ok := expr.(*sqlparser.AliasedExpr)
		if !ok {
			return nil, fmt.Errorf("elasticsql: unsupported star expression in function date_range")
		}

		switch item := nonStarExpr.Expr.(type) {
		case *sqlparser.ComparisonExpr:
			colName := sqlparser.String(item.Left)
			equalVal := sqlparser.String(item.Right.(*sqlparser.SQLVal))
			// fmt.Printf("%#v", sqlparser.String(item.Right))
			equalVal = strings.Trim(equalVal, `'`)

			switch colName {
			case "field":
				field = equalVal
			case "format":
				format = equalVal
			default:
				return nil, fmt.Errorf("elasticsql: unsupported column name " + colName)
			}
		case *sqlparser.SQLVal:
			skippedString := strings.Trim(sqlparser.String(item), "`")
			rangeList = append(rangeList, skippedString)
		default:
			return nil, fmt.Errorf("elasticsql: unsupported expression " + sqlparser.String(expr))
		}
	}

	if len(field) == 0 {
		return nil, fmt.Errorf("elasticsql: lack field of date_range")
	}

	for i := 0; i < len(rangeList)-1; i++ {
		tmpMap := msi{
			"from": strings.Trim(rangeList[i], `'`),
			"to":   strings.Trim(rangeList[i+1], `'`),
		}
		rangeMapList = append(rangeMapList, tmpMap)
	}

	innerMap = msi{
		"date_range": msi{
			"field":  field,
			"ranges": rangeMapList,
			"format": format,
		},
	}

	return innerMap, nil
}

func handleGroupByFuncExpr(funcExpr *sqlparser.FuncExpr, child msi) (msi, error) {
	var innerMap msi
	var err error

	switch funcExpr.Name.Lowered() {
	case "date_histogram":
		innerMap, err = handleGroupByFuncExprDateHisto(funcExpr)
	case "range":
		innerMap, err = handleGroupByFuncExprRange(funcExpr)
	case "date_range":
		innerMap, err = handleGroupByFuncExprDateRange(funcExpr)
	default:
		return nil, fmt.Errorf("elasticsql: unsupported group by functions" + sqlparser.String(funcExpr))
	}

	if err != nil {
		return nil, err
	}

	if len(child) > 0 && innerMap != nil {
		innerMap["aggregations"] = child
	}

	stripedFuncExpr := sqlparser.String(funcExpr)
	stripedFuncExpr = strings.Replace(stripedFuncExpr, " ", "", -1)
	stripedFuncExpr = strings.Replace(stripedFuncExpr, "'", "", -1)
	return msi{stripedFuncExpr: innerMap}, nil
}

func handleGroupByAgg(groupBy sqlparser.GroupBy, innerMap msi) (msi, error) {
	aggMap := make(msi)

	child := innerMap

	for i := len(groupBy) - 1; i >= 0; i-- {
		v := groupBy[i]

		switch item := v.(type) {
		case *sqlparser.ColName:
			currentMap := handleGroupByColName(item, i, child)
			child = currentMap

		case *sqlparser.FuncExpr:
			currentMap, err := handleGroupByFuncExpr(item, child)
			if err != nil {
				return nil, err
			}
			child = currentMap
		}
	}
	aggMap = child

	return aggMap, nil
}

func buildAggs(sel *sqlparser.Select) (string, error) {
	funcExprArr, _, funcErr := extractFuncAndColFromSelect(sel.SelectExprs)
	innerAggMap := handleFuncInSelectAgg(funcExprArr)

	if funcErr != nil {
	}

	aggMap, err := handleGroupByAgg(sel.GroupBy, innerAggMap)
	if err != nil {
		return "", err
	}

	mapJSON, _ := json.Marshal(aggMap)

	return string(mapJSON), nil
}

// extract func expressions from select exprs
func extractFuncAndColFromSelect(sqlSelect sqlparser.SelectExprs) ([]*sqlparser.FuncExpr, []*sqlparser.ColName, error) {
	var colArr []*sqlparser.ColName
	var funcArr []*sqlparser.FuncExpr
	for _, v := range sqlSelect {
		// non star expressioin means column name
		// or some aggregation functions
		expr, ok := v.(*sqlparser.AliasedExpr)
		if !ok {
			// no need to handle, star expression * just skip is ok
			continue
		}

		// NonStarExpr start
		switch expr.Expr.(type) {
		case *sqlparser.FuncExpr:
			funcExpr := expr.Expr.(*sqlparser.FuncExpr)
			funcArr = append(funcArr, funcExpr)

		case *sqlparser.ColName:
			continue
		default:
			// ignore
		}

		// starExpression like *, table.* should be ignored
		// 'cause it is meaningless to set fields in elasticsearch aggs
	}
	return funcArr, colArr, nil
}
