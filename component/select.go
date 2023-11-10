package component

import (
	"reflect"
	"strings"
)

type SelectType string

const (
	Input       SelectType = "input"
	Sel         SelectType = "select"
	InputSelect SelectType = "input_select"
	BoolSelect  SelectType = "bool_select"
	SingleDate  SelectType = "single_date"
	BetweenDate SelectType = "between_date"
)

type Select struct {
	SelectType  SelectType `json:"select_type"`
	ApiPath     string     `json:"api_path,omitempty"`
	Column      string     `json:"column,omitempty"`
	Placeholder string     `json:"placeholder"`
}

const (
	tag           = "select"
	defaultColTag = "dc"
)

func ToSelectRes(data any) any {
	var res = make(map[string]any)
	objT := reflect.TypeOf(data)
	objV := reflect.ValueOf(data)
	if objT.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < objT.NumField(); i++ {
		tagV := objT.Field(i).Tag.Get(tag)
		if tagV == "" {
			continue
		} else {
			res[tagV] = objV.Field(i).Interface()
		}

	}
	return res
}

// MapToSelectRes json 转 select
func MapToSelectRes(data []map[string]any, model any) []map[string]any {

	var res = make([]map[string]any, 0, len(data))
	var mapping = make(map[string]string)
	objT := reflect.TypeOf(model)
	if objT.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < objT.NumField(); i++ {
		tagV := objT.Field(i).Tag.Get(tag)
		json := strings.Split(objT.Field(i).Tag.Get(jsonTag), jsonSplit)[0]
		if json == "" {
			continue
		}
		if tagV == "" {
			continue
		} else {
			mapping[json] = tagV
		}

	}
	for _, datum := range data {
		var re = make(map[string]any)
		for k, v := range datum {
			if newKey, exist := mapping[k]; exist {
				re[newKey] = v
			}
		}
		res = append(res, re)
	}
	return res
}

func ToDefaultCol(data any) []string {
	var res []string
	objT := reflect.TypeOf(data)
	if objT.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < objT.NumField(); i++ {
		json := strings.Split(objT.Field(i).Tag.Get(jsonTag), jsonSplit)[0]
		if json == "" {
			continue
		}
		_, ok := objT.Field(i).Tag.Lookup(defaultColTag)
		if ok {
			res = append(res, json)
		}
	}
	return res
}
