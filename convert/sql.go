package convert

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

const andTag = "and"

var nullMap = map[any]struct{}{
	"无":         {},
	"null":      {},
	"nil":       {},
	"undefined": {},
	"-":         {},
}

// AndSql FIXME 1.防sql注入  2.优化圈复杂度  3.优化代码逻辑层次
//
//	用于将任意结构体转 and sql
//	tag: `and`
func AndSql(data any) string {
	if data == nil {
		return ""
	}
	if reflect.ValueOf(data).IsZero() {
		return ""
	}
	var res = ""
	var reqMap, omitKey = map[string]any{}, map[string]struct{}{}
	var likeMap = make(map[string]struct{}, 0)
	tagT, tagV := reflect.TypeOf(data), reflect.ValueOf(data)
	for i := 0; i < tagT.NumField(); i++ {
		now := tagT.Field(i)
		jsonKey := strings.Split(now.Tag.Get("json"), ",")[0]
		if tagV.Field(i).IsZero() || tagV.Field(i).Len() == 0 {
			omitKey[jsonKey] = struct{}{}
		}
		if v, ok := now.Tag.Lookup(andTag); ok {
			switch v {
			case "like":
				likeMap[jsonKey] = struct{}{}
			case "-":
				omitKey[jsonKey] = struct{}{}
			}
		}
	}
	marshal, err := json.Marshal(data)
	if err != nil || len(marshal) == 0 || string(marshal) == "" {
		return ""
	}

	err = json.Unmarshal(marshal, &reqMap)
	if err != nil {
		return ""
	}
	for key, value := range reqMap {
		if _, ok := omitKey[key]; ok {
			continue
		}
		objV, objT := reflect.ValueOf(value), reflect.TypeOf(value)
		if objT.Kind() < reflect.Array || objT.Kind() == reflect.String {
			if _, ok := nullMap[value]; ok {
				res += fmt.Sprintf("and %s is null ", key)
				continue
			}
			if !objV.IsZero() {
				if _, ok := likeMap[key]; ok {
					res += fmt.Sprintf("and %s like '%%%s%%' ", key, TransSql(value.(string)))
				} else {
					res += fmt.Sprintf("and %s = '%s' ", key, TransSql(value.(string)))
				}
			}
		} else if !objV.IsNil() {
			if objT.Kind() == reflect.Slice || objT.Kind() == reflect.Array {
				vs := value.([]any)
				// 如果是字符串数组,且长度为1,按逗号分隔
				if len(vs) == 1 {
					if str, ok := vs[0].(string); ok {
						var a []any
						for _, s := range strings.Split(str, ",") {
							a = append(a, s)
						}
						vs = a
					}
				}

				if _, ok := nullMap[vs[0]]; ok {
					res += fmt.Sprintf("and %s is null ", key)
					continue
				}
				tem, symbol, base := "", ",", "and  ( %s in (%s) %s"
				_, ok := likeMap[key]
				if ok {
					symbol = "|"
					base = "and (%s regexp "
				}
				var flag bool
				for i, v := range vs {
					_, flag = nullMap[v]
					if flag {
						continue
					}
					bytes, _ := json.Marshal(v)
					if !ok {
						tem += fmt.Sprintf("%s", string(bytes))
					} else {
						tem += string(bytes)
					}
					tem += symbol
					if i == len(vs)-1 {
						tem = strings.TrimRight(tem, symbol)
					}
				}
				var last = ")"
				if flag {
					last = fmt.Sprintf(` or %s is null )`, key)
				}
				res += fmt.Sprintf(base, key, tem, last)

			}
		}
	}
	return res
}

func GormAndSql(db *gorm.DB, data any) *gorm.DB {

	return db.Where(" 1=1 " + AndSql(data))
}

func SortBy(sortList ...string) string {
	var res = ` order by`
	for i, s := range sortList {
		if strings.Trim(s, " ") == "" {
			continue
		}
		res += " " + s
		if i != len(sortList)-1 {
			res += `,`
		}
	}
	res = strings.TrimRight(res, " ")
	if strings.HasSuffix(res, `order by`) {
		return " "
	}
	return res + "  "
}

// TransSql 简单的转义字符串防注入
func TransSql(sql string) string {
	return strings.ReplaceAll(strings.ReplaceAll(sql, `'`, `\'`), `"`, `\"`)
}

// AndBetweenSql Between
func AndBetweenSql(col string, data []string, notInclude bool) string {
	if len(data) == 1 {
		data = strings.Split(data[0], ",")
	}
	var operator = "between"
	if len(data) == 1 {
		return fmt.Sprintf(` and %s >= '%s'`, col, data[0])
	} else if len(data) > 1 {
		if notInclude {
			operator = " not between "
		}
		return fmt.Sprintf(` and ( %s %s '%s' and '%s') `, col, operator, data[0], data[1])
	}
	return ""
}
