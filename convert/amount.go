package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Country string

const (
	AmountMatch = "^(-|−| )?(\\d{1,3}((.|,| | | )\\d{3})*|\\d+)((\\.|,)\\d{2})?$"
)

func ToNormalAmount(str string) string {
	if _, err := strconv.Atoi(str); err == nil {
		return str
	}

	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return fmt.Sprintf(`%.2f`, f)
	}

	return fmt.Sprintf(`%.2f`, NewAmount(str).SplitAll(".", ",", " ", " ", " ").Float64())
}

type amount struct {
	isNegative bool
	str        string
	slice      []string
}

func IsVailAmount(str string) bool {
	if ok, _ := regexp.MatchString(AmountMatch, str); ok {
		return true
	}
	return false
}

func NewAmount(str string) *amount {
	str = strings.ReplaceAll(str, "−", "-")

	var isNegative bool
	cutNegative := strings.TrimPrefix(str, `-−`)
	if len(cutNegative) != len(str) {
		isNegative = true
		str = cutNegative
	}
	return &amount{
		isNegative: isNegative,
		str:        str,
		slice:      []string{},
	}
}

func (a *amount) TrimLeft(cutest string) *amount {
	a.str = strings.TrimLeft(a.str, cutest)
	return a
}

func (a *amount) ReplaceAll(dst string, cuts ...string) *amount {
	for _, cut := range cuts {
		a.str = strings.ReplaceAll(a.str, cut, dst)
	}
	return a
}

func (a *amount) String() string {
	return a.str
}

func (a *amount) SplitAll(points ...string) *amount {
	a.slice = []string{a.str}
	for _, point := range points {
		var temp []string
		for _, s := range a.slice {
			temp = append(temp, strings.Split(s, point)...)
		}
		a.slice = temp
	}

	return a
}

func (a *amount) Float64() float64 {
	var fraction string
	// 最后一个长度为 2 就认为是小数点后的
	lastIndex := len(a.slice) - 1
	last := a.slice[lastIndex]
	if len(last) == 2 {
		fraction = last
	}
	var str string
	if fraction == "" {
		str = strings.Join(a.slice, "")
	} else {
		str = strings.Join(a.slice[:lastIndex], "") + "." + fraction
	}
	float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return float

}

func (a *amount) ParseFloat64() float64 {
	return a.SplitAll(".", ",", " ").Float64()
}
