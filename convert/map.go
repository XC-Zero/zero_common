package convert

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	HighLight      = "hl"
	HighLightGreen = "hl_green"
	HighLightBlue  = "hl_blue"
)

type PrimaryMap map[string]map[string]any

func ToPrimaryMap(json []map[string]any, keys ...string) PrimaryMap {
	var res = make(PrimaryMap, len(json))

	for _, m := range json {
		key := ""
		for _, o := range keys {
			v, _ := m[o]
			key += fmt.Sprintf("%v", v)
		}
		res[key] = m
	}
	return res
}

// Different 比较两个map的区别
func Different(a, b PrimaryMap, colMapping map[string]string) PrimaryMap {
	var lock sync.Mutex
	id := uuid.New().String()
	log.Println("Start to check different! ID is ", id)
	var res = make(PrimaryMap, (len(a)+len(b))>>1)
	intersection := GetIntersectionKey(a, b)
	log.Println("intersection len is ", len(intersection))
	wg := sync.WaitGroup{}

	wg.Add(3)

	go func(inter map[string]struct{}) {
		for mKey, mObj := range a {
			if _, ok := inter[mKey]; !ok {
				lock.Lock()
				res[AddHtml(mKey, HighLightGreen)] = map[string]any{"a": mObj}
				lock.Unlock()
			}
		}
		wg.Done()
	}(intersection)

	go func(inter map[string]struct{}) {
		for mKey, mObj := range b {
			if _, ok := inter[mKey]; !ok {
				lock.Lock()
				res[AddHtml(mKey, HighLightBlue)] = map[string]any{"b": mObj}
				lock.Unlock()
			}
		}
		wg.Done()
	}(intersection)

	// 核对交集里的
	go func() {
		for key := range intersection {
			var tempA, tempB = make(map[string]any, 0), make(map[string]any, 0)
			var flag = false
			for ak, av := range a[key] {
				if bk, ok := colMapping[ak]; ok {
					bv := b[key][bk]
					if bv != av {
						flag = true
						tempA[AddHtml(ak, HighLightBlue)] = av
						tempB[AddHtml(ak, HighLightGreen)] = bv
					}
				}

			}
			if flag == true {
				lock.Lock()
				res[key] = map[string]any{
					"a": tempA,
					"b": tempB,
				}
				lock.Unlock()
			}
		}
		wg.Done()

	}()

	wg.Wait()
	log.Println("End to check different! ID is ", id)

	return res
}

// GetIntersectionKey 获取两个map的交集
func GetIntersectionKey[T any | map[string]string](a map[string]T, b map[string]T) map[string]struct{} {
	var res = make(map[string]struct{}, 0)
	log.Println(len(a), len(b))
	var empty = struct{}{}
	for s := range a {
		if _, ok := b[s]; ok {
			res[s] = empty
		}
	}
	for s := range b {
		if _, ok := a[s]; ok {
			res[s] = empty
		}
	}
	return res

}

func GetUnionKey[T any | map[string]string](a map[string]T, b map[string]T) map[string]struct{} {
	var res = make(map[string]struct{}, len(a))
	var empty = struct{}{}

	for s := range a {
		res[s] = empty
	}
	for s := range b {
		res[s] = empty
	}
	return res
}

func GetDifferentKey[T any | map[string]string](a map[string]T, b map[string]T) map[string]struct{} {
	var res = make(map[string]struct{}, 0)
	var empty = struct{}{}

	var inter = GetIntersectionKey(a, b)
	for s := range GetUnionKey(a, b) {
		if _, ok := inter[s]; !ok {
			res[s] = empty
		}
	}
	return res
}

func AddHtml(content string, htmlTag string) string {
	return "<" + htmlTag + ">" + content + "</" + htmlTag + ">"
}

// MapToGormStruct json转gorm标签的结构体
func MapToGormStruct(data map[string]any, structPtr any) error {
	//spew.Dump(structPtr)
	if reflect.TypeOf(structPtr).Kind() != reflect.Ptr {
		return errors.New("Type of structPtr is not a pointer!   ")
	}
	objV := reflect.ValueOf(structPtr).Elem()
	if reflect.TypeOf(objV).Kind() != reflect.Struct {
		return errors.New("Type of structPtr is not a struct pointer! It is " + reflect.TypeOf(objV).String())
	}
	objT := reflect.TypeOf(structPtr).Elem()

	for i := 0; i < objT.NumField(); i++ {
		// 匿名嵌套 或 不可访问的 跳过
		if objT.Field(i).Anonymous || !IsFirstUpper(objT.Field(i).Name) {
			continue
		}
		tags := strings.Split(objT.Field(i).Tag.Get(`gorm`), ";")
		var colName string
		for _, tag := range tags {
			if strings.HasPrefix(tag, "column:") {
				colName = tag[7:]
				break
			}
		}
		if newVal, exist := data[colName]; exist {

			newV := reflect.ValueOf(newVal)
			tarType := objT.Field(i).Type
			if newV.IsValid() && !newV.IsZero() {
				souType := reflect.TypeOf(newVal)

				// 如果可以直接转则直接转
				if newV.CanConvert(tarType) {
					objV.Field(i).Set(newV.Convert(tarType))
					//	 指针转非指针
				} else if souType.Kind() == reflect.Ptr {
					if reflect.TypeOf(souType.Elem()) == tarType {
						objV.Field(i).Set(newV.Elem())
					} else {
						if tarType.Elem().Kind() == reflect.Int {
							switch newVal.(type) {
							case float64:
								var a = newVal.(float64)
								objV.Field(i).Set(reflect.ValueOf(int(a)))
								continue
							case string:
								var a = newVal.(string)
								aa, err := strconv.Atoi(a)
								if err == nil {
									objV.Field(i).Set(reflect.ValueOf(aa))
								}
							}
						}

					}
					//	 非指针转指针
				} else if tarType.Kind() == reflect.Ptr {
					//	 类型相同则转个指针
					if souType == tarType.Elem() {
						switch newVal.(type) {
						case string:
							var a = newVal.(string)
							objV.Field(i).Set(reflect.ValueOf(&a))
							continue
						case int:
							var a = newVal.(int)
							objV.Field(i).Set(reflect.ValueOf(&a))
							continue
						case int32:
							var a = newVal.(int32)
							objV.Field(i).Set(reflect.ValueOf(&a))
							continue
						case int64:
							var a = newVal.(int64)
							objV.Field(i).Set(reflect.ValueOf(&a))
							continue
						case float64:
							var a = newVal.(float64)
							objV.Field(i).Set(reflect.ValueOf(&a))
							continue
						case bool:
							var a = newVal.(bool)
							objV.Field(i).Set(reflect.ValueOf(&a))
							continue
						}

					} else {
						if tarType.Elem().Kind() == reflect.Int {
							switch newVal.(type) {
							case float64:
								var a = newVal.(float64)
								var aa = int(a)
								objV.Field(i).Set(reflect.ValueOf(&aa))
								continue
							case string:
								var a = newVal.(string)
								aa, err := strconv.Atoi(a)
								if err == nil {
									objV.Field(i).Set(reflect.ValueOf(&aa))
								}
							}
						}
						if tarType.Elem().String() == "time.Time" {
							switch newVal.(type) {
							case int64, int, int32:
								var a = newVal.(int64)
								var aa = time.Unix(a/1000, 0).Add(-time.Hour * 8)
								objV.Field(i).Set(reflect.ValueOf(&aa))
								continue
							case float64, float32:
								var a = newVal.(float64)
								var aa = time.Unix(int64(a)/1000, 0).Add(-time.Hour * 8)
								objV.Field(i).Set(reflect.ValueOf(&aa))
								continue
							case string:
								var a = newVal.(string)
								parse, _ := time.Parse(time.RFC3339, a)
								objV.Field(i).Set(reflect.ValueOf(&parse))
							}
						}
					}
				} else if souType.Kind() == reflect.String && tarType.String() != "time.Time" {
					s := newVal.(string)
					switch tarType.Kind() {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						atoi, err := strconv.Atoi(s)
						if err == nil {
							objV.Field(i).Set(reflect.ValueOf(atoi).Convert(tarType))
						}
					case reflect.Float64, reflect.Float32:
						b, err := strconv.ParseFloat(s, 64)
						if err == nil {
							objV.Field(i).Set(reflect.ValueOf(b).Convert(tarType))
						}
					case reflect.Bool:
						b, err := strconv.ParseBool(s)
						if err == nil {
							objV.Field(i).Set(reflect.ValueOf(b))
						}
					}
					if tarType.String() == "time.Time" {
						var a = newVal.(string)
						parse, err := time.Parse(time.RFC3339, a)
						objV.Field(i).Set(reflect.ValueOf(parse))
						if err != nil {
							panic(err)
						}
					}
				} else if tarType.String() == "time.Time" {
					switch newVal.(type) {
					case int64, int, int32:
						var a = newVal.(int64)
						var aa = time.Unix(a/1000, 0).Add(-time.Hour * 8)
						objV.Field(i).Set(reflect.ValueOf(aa))
						continue
					case float64, float32:
						var a = newVal.(float64)
						var aa = time.Unix(int64(a)/1000, 0).Add(-time.Hour * 8)
						objV.Field(i).Set(reflect.ValueOf(aa))
						continue
					case string:
						var a = newVal.(string)
						parse, _ := time.Parse(time.RFC3339, a)
						objV.Field(i).Set(reflect.ValueOf(parse))

					}
				}
			}
		}

	}
	return nil
}

func MapMapping(data map[string]any, mapping map[string]string) map[string]any {
	var res = make(map[string]any)
	for org, tar := range mapping {
		if v, exist := data[org]; exist {
			res[tar] = v
		}
	}
	return res
}

func StructToMap(structPtr any) (res map[string]any, err error) {
	if reflect.TypeOf(structPtr).Kind() != reflect.Ptr {
		err = errors.New("Type of structPtr is not a pointer!   ")
		return
	}
	objV := reflect.ValueOf(structPtr).Elem()
	if reflect.TypeOf(objV).Kind() != reflect.Struct {
		err = errors.New("Type of structPtr is not a struct pointer! It is " + reflect.TypeOf(objV).String())
		return
	}
	marshal, err := json.Marshal(structPtr)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		return nil, err
	}
	for key, value := range res {
		if strings.Contains(key, ";") {
			res[strings.Split(key, ";")[0]] = value
			delete(res, key)
		}
		if len(key) == 0 {
			delete(res, key)
		}
	}
	return
}

func SliceToMap(a any, omitKey ...string) ([]map[string]any, error) {
	bb, err := json.Marshal(a)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var mm []map[string]any
	err = json.Unmarshal(bb, &mm)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for i, mi := range mm {
		for _, key := range omitKey {
			delete(mi, key)
		}
		mm[i] = mi
	}

	return mm, nil
}

// 结构体转元组
func StructToTuple(str any, prefix, suffix string) string {
	var s []string

	objV := reflect.ValueOf(str)
	for i := 0; i < objV.NumField(); i++ {
		var vv = objV.Field(i).Interface()
		var strin string
		switch vv.(type) {
		case string:
			strin = vv.(string)
		case int, int32, int64, int16, int8, uint:
			strin = fmt.Sprintf(`%d`, vv)
		case float64, float32:
			strin = fmt.Sprintf(`%f`, vv)
		}
		s = append(s, prefix+strin+suffix)
	}
	return `(` + strings.Join(s, ",") + `)`
}
