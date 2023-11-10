package convert

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strings"
)

// MapperStructByName
//
//	根据字段名转换 仅支持 结构体指针
func MapperStructByName(dst, src any) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("Type of dst is not ptr!")
	}
	if reflect.TypeOf(src).Kind() != reflect.Ptr {
		return errors.New("Type of src is not ptr!")
	}

	from := reflect.ValueOf(src).Elem()
	fromTyp := reflect.TypeOf(src).Elem()
	to := reflect.ValueOf(dst).Elem()
	toTyp := reflect.TypeOf(dst).Elem()
	if from.Kind() != reflect.Struct {
		return errors.New("Type of src is not struct ptr!")
	}
	if to.Kind() != reflect.Struct {
		return errors.New("Type of dst is not struct ptr!")
	}

	for i := 0; i < toTyp.NumField(); i++ {
		firstLetter := strings.Split(toTyp.Field(i).Name, "")[0]
		if strings.ToLower(firstLetter) == firstLetter {
			continue
		}
		srcName := toTyp.Field(i).Name
		v := from.FieldByName(srcName)
		vt, _ := fromTyp.FieldByName(srcName)
		if v.IsValid() && !v.IsZero() {
			gormTag := vt.Tag.Get("gorm")
			if strings.Contains(gormTag, `type:decimal(`) {

				if v.Type().Kind() == reflect.String {
					s, _ := v.Interface().(string)
					v = reflect.ValueOf(TrimRightZeroRetain(s, 2))

				}
				if v.Type().Kind() == reflect.Ptr && v.Elem().Type().Kind() == reflect.String {
					s, _ := v.Interface().(*string)
					if s != nil {
						t := TrimRightZeroRetain(*s, 2)
						v = reflect.ValueOf(&t)
					}
				}
			}
			to.FieldByName(srcName).Set(v)

		}

	}
	return nil

}

func CompareDifference(a, b any, aFlag, bFlag, tag string, retainKeys ...string) (map[string]any, error) {

	aT, bT := reflect.TypeOf(a), reflect.TypeOf(b)
	if aT != bT {
		return nil, errors.Errorf(`Difference type a and b, a is %s , b is %s `, aT, bT)
	}

	retainMap := make(map[string]struct{})
	for _, key := range retainKeys {
		retainMap[key] = struct{}{}
	}
	if aFlag == "" {
		aFlag = "a"
	}
	if bFlag == "" {
		bFlag = "b"
	}

	res := make(map[string]any)
	aV, bV := reflect.ValueOf(a), reflect.ValueOf(b)

	for i := 0; i < aT.NumField(); i++ {
		now := aT.Field(i)
		if !now.Anonymous {
			name := now.Name
			tagName := name
			if !IsFirstUpper(name) {
				continue
			}
			if tag != "" {
				tn := strings.Split(now.Tag.Get(tag), ",")[0]
				if tn != "" {
					tagName = tn
				}
			}

			nowAV, nowBV := aV.FieldByName(name), bV.FieldByName(name)
			if now.Type.Kind() == reflect.Ptr {
				if !nowAV.IsZero() {
					nowAV = nowAV.Elem()
				}
				if !nowBV.IsZero() {
					nowBV = nowBV.Elem()
				}
			}
			if tagName == `delivery_plan_time` {
				spew.Dump(nowAV.Interface(), nowBV.Interface())
			}
			//如果不相等
			if !nowAV.Equal(nowBV) {
				aa := make(map[string]any)
				aa[aFlag], aa[bFlag] = nowAV.Interface(), nowBV.Interface()
				res[tagName] = aa
			} else if _, exist := retainMap[name]; exist {
				res[tagName] = nowAV.Interface()
			}

		}
	}
	return res, nil
}

// InvisibleDataCol
//
//	从 结构体,结构体切片,map 中将指定字段(允许json标签)置为 空值
//	不支持嵌套结构体
func InvisibleDataCol(data any, columns ...string) any {
	dataType := reflect.TypeOf(data)
	switch dataType.Kind() {
	case reflect.Map:
		dataVal := reflect.ValueOf(data)
		for i := range columns {
			col := reflect.ValueOf(columns[i])
			vv := dataVal.MapIndex(col)
			if vv.IsValid() && !vv.IsNil() && !vv.IsZero() {
				dataVal.SetMapIndex(col, reflect.Zero(reflect.TypeOf(vv)))
			}
		}
		return data
	case reflect.Struct:
		vv := reflect.ValueOf(data)
		tagMap := StructTagMapIndex(data, "json")
		for i := 0; i < len(columns); i++ {
			vvv := vv.FieldByName(columns[i])
			if vvv.IsValid() && !vvv.IsNil() && !vvv.IsZero() {
				log.Printf("set column %s zero", columns[i])
				vvv.SetZero()
			} else {
				if idx, exist := tagMap[columns[i]]; exist {
					nowField := vv.FieldByIndex(idx)
					log.Println(nowField.String())
					if nowField.CanSet() {
						log.Printf("set column %s zero", columns[i])
						vv.FieldByIndex(idx).SetZero()
					} else {
						log.Printf(`[WARN] This column %s cannot be set`, columns[i])
					}
				}
			}

		}

		return vv.Interface()

	case reflect.Ptr:
		return InvisibleDataCol(reflect.ValueOf(data).Elem().Interface(), columns...)

	case reflect.Slice, reflect.Array:
		val := reflect.ValueOf(data)
		tagMap := StructTagMapIndex(data, "json")
		// FIXME 当前只支持 结构体切片 需要支持 map切片
		for i := 0; i < val.Len(); i++ {
			for _, column := range columns {
				if idx, exist := tagMap[column]; exist {
					val.Index(i).FieldByIndex(idx).SetZero()
				}
			}
		}
		return data
	default:
		return data
	}

}

func StructTagMapIndex(structVal any, tag string) map[string][]int {
	typ := reflect.TypeOf(structVal)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	var tagMap = make(map[string][]int)

	for j := 0; j < typ.NumField(); j++ {
		jsonTag := typ.Field(j).Tag.Get(tag)
		if jsonTag == "" {
			continue
		} else {
			tagMap[strings.Split(jsonTag, ",")[0]] = typ.Field(j).Index
		}
	}
	return tagMap
}

//func InvisibleMap(any2 any)  {
//
//}
