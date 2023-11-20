package component

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type Table struct {
	HasStripe         bool     `json:"has_stripe"`
	HasChildrenTable  bool     `json:"has_children_table"`
	HasCheck          bool     `json:"has_check"`
	HasBatchEdit      bool     `json:"has_batch_edit"`
	TableColumns      []Col    `json:"table_columns"`
	ChildTableColumns []Col    `json:"child_table_columns"`
	TableSelects      []Select `json:"table_selects"`
	TableBatchEdit    []Select `json:"table_batch_edit"`
}

type Detail struct {
	DetailName string `json:"detail_name"`
	DetailCols []Col  `json:"detail_cols"`
	CanEdit    bool   `json:"can_edit "`
}

type FixCol string

const (
	Left  FixCol = "left"
	Right FixCol = "right"
	NIL   FixCol = ""
)

type ColValType string

const (
	STRING    ColValType = "string"
	DECIMAL   ColValType = "decimal"
	BOOL      ColValType = "bool"
	INT       ColValType = "int"
	TEXT      ColValType = "text"
	DATETIME  ColValType = "datetime"
	FLOAT     ColValType = "float"
	OPERATION ColValType = "operation"
	UNKNOW    ColValType = ""
)

var OperationCol = Col{
	ColLabel:   "操作",
	ColVal:     "operation",
	ColValType: OPERATION,
}

type Col struct {
	ColLabel     string     `json:"title"`
	ColVal       string     `json:"dataIndex"`
	ColValType   ColValType `json:"col_val_type"`
	CanEdit      bool       `json:"can_edit"`
	CanBatchEdit bool       `json:"can_batch_edit"`
	CanSort      bool       `json:"can_sort"`
	CanHover     bool       `json:"can_hover"`
	CanClick     bool       `json:"can_click"`
	Resizable    bool       `json:"resizable"`
	FixCol       FixCol     `json:"fixed"`
	Width        *int       `json:"width"`
}

const (
	tableTag       = "table"
	tableSplit     = ";"
	tableLabel     = "label"
	tableCanSort   = "sort"
	tableCanHover  = "hover"
	tableCanClick  = "click"
	edit           = "edit"
	tableResizable = "resizable"
	tableEqu       = ":"
	tableFix       = "fixed"
	jsonTag        = "json"
	jsonSplit      = ","
)

func NewTable(structVal any) []Col {
	var cols []Col

	objT := reflect.TypeOf(structVal)

	if objT.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < objT.NumField(); i++ {
		var col Col
		json := strings.Split(objT.Field(i).Tag.Get(jsonTag), jsonSplit)[0]
		if json == "" {
			continue
		}
		col.ColVal = json
		var tagMap = make(map[string]string)
		tags := objT.Field(i).Tag.Get(tableTag)
		kind := objT.Field(i).Type.Kind()
		if kind == reflect.Ptr {
			kind = objT.Field(i).Type.Elem().Kind()
		}
		col.ColValType = ColValType(kind.String())

		for _, tag := range strings.Split(tags, tableSplit) {
			tagContent := strings.Split(tag, tableEqu)
			if len(tagContent) != 2 {
				continue
			}
			tagMap[tagContent[0]] = tagContent[1]
		}
		if v, ok := tagMap[tableLabel]; ok {
			if v == "-" {
				continue
			}
			col.ColLabel = v
		} else {
			col.ColLabel = json
		}
		if v, ok := tagMap[tableCanSort]; ok {
			col.CanSort = v == "true"
		}
		if v, ok := tagMap[tableFix]; ok {
			col.FixCol = FixCol(fmt.Sprintf("%v", v))
		}
		if v, ok := tagMap[tableCanClick]; ok {
			col.CanClick = v == "true"
		}
		if v, ok := tagMap[tableCanHover]; ok {
			col.CanHover = v == "true"
		}
		if v, ok := tagMap[edit]; ok {
			col.CanEdit = v == "true" || v == "batch"
			col.CanBatchEdit = v == "batch"
		}
		if v, ok := tagMap[tableResizable]; ok && v == "false" {
			col.Resizable = false
		} else {
			col.Resizable = true
		}
		cols = append(cols, col)
	}

	return cols
}

// ColsToMap 返回列的  英 -> 中文 对照
func ColsToMap(cols []Col) map[string]string {
	var mm = make(map[string]string, len(cols))
	for _, col := range cols {
		mm[col.ColVal] = col.ColLabel
	}
	return mm
}

// ReverseColsToMap 返回列的  中 -> 英文对照
func ReverseColsToMap(cols []Col) map[string]string {
	var mm = make(map[string]string, len(cols))
	for _, col := range cols {
		mm[col.ColLabel] = col.ColVal
	}
	return mm
}

// MapToColsMap 实体的JSON转中英文
func MapToColsMap(data []map[string]any, model any) ([]Col, []map[string]any) {
	var res = make([]map[string]any, 0, len(data))
	cols := NewTable(model)
	mapping := ColsToMap(cols)
	for _, datum := range data {
		var mm = make(map[string]any, len(datum))
		for k, v := range datum {
			if newKey, exist := mapping[k]; exist {
				mm[newKey] = v
			}
		}
		if len(mm) == 0 {
			continue
		}
		res = append(res, mm)
	}
	return cols, res
}

// StructToColMap 结构体数组转json数组
func StructToColMap(structSlice any) (res []map[string]any, err error) {
	var mapping map[string]string
	switch reflect.TypeOf(structSlice).Kind() {
	case reflect.Slice, reflect.Array:
		refVal := reflect.ValueOf(structSlice)
		if refVal.Len() == 0 {
			return
		}
		first := refVal.Index(0)
		if first.Type().Kind() != reflect.Struct {
			return nil, errors.WithStack(errors.New("Unsupported slice type! Type is :" + first.Type().Kind().String()))
		}
		mapping = ColsToMap(NewTable(first.Interface()))
	case reflect.Struct:
		mapping = ColsToMap(NewTable(structSlice))
	default:
		return nil, errors.WithStack(errors.New("Unsupported type!"))
	}
	marshal, err := json.Marshal(structSlice)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i, re := range res {
		for key, value := range re {
			if newKey, exist := mapping[key]; exist {
				delete(re, key)
				re[newKey] = value
			}
		}
		res[i] = re
	}
	return res, nil
}

type TableHookFunc func(any2 any, table Table) (any, Table)
