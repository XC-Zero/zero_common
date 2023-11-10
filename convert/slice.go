package convert

import (
	"strconv"
	"strings"
)

func StrSliceToIntSlice(a []string) (b []int) {
	for _, s := range a {
		i, _ := strconv.Atoi(s)

		b = append(b, i)
	}
	return
}

func StrToIntSlice(str, symbol string) []int {

	return StrSliceToIntSlice(strings.Split(str, symbol))
}

func IntOrInt(a, b []int) []int {
	for i := 0; i < len(a); i++ {
		a[i] = a[i] | b[i]
	}
	return a
}

func IntJoin(a []int, symbol string) string {
	var res string
	for i := 0; i < len(a); i++ {
		res += strconv.Itoa(a[i])
		if i != len(a)-1 {
			res += symbol
		}
	}
	return res
}

func StringToIntSlice(str string, symbol string) (res []int) {

	s := strings.Split(str, symbol)
	for _, s2 := range s {
		in, err := strconv.Atoi(s2)
		if err == nil {
			res = append(res, in)
		}
	}
	if len(res) == 0 {
		res = nil
	}
	return
}

func SliceContain[T int | int32 | int64 | float64 | byte | string](list []T, str T) bool {
	for i := range list {
		if list[i] == str {
			return true
		}
	}
	return false
}

func GetValue[T any](list []T, key int) T {
	var d T
	if len(list) > key {
		return list[key]
	}
	return d
}
