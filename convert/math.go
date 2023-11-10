package convert

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var chineseMapping = map[string]string{
	"0": "零",
	"1": "壹",
	"2": "贰",
	"3": "弎",
	"4": "肆",
	"5": "伍",
	"6": "陆",
	"7": "柒",
	"8": "捌",
	"9": "玖",
}

var chineseUnit = map[int]string{
	1:  "毫",
	2:  "厘",
	3:  "分",
	4:  "角",
	7:  "拾",
	8:  "佰",
	9:  "仟",
	10: "万",
	13: "亿",
}

const (
	tenThousand    = 100 * 100
	hundredMillion = 10000 * 10000
)

// ToUpperChineseDecimal 转中文大写的小数
func ToUpperChineseDecimal(f float64) string {
	var res = ""
	integer := int(f)
	res += ToUpperChinese(integer) + "元"
	fraction := strings.Split(fmt.Sprintf("%.2f", f), ".")[1]
	data := strings.Split(fraction, "")
	if fraction == "00" {
		res += "整"
	} else {
		j := chineseMapping[data[0]]
		res += j
		if data[1] != "0" {
			if data[0] != "0" {
				res += "角"
			}
			res += chineseMapping[data[1]] + "分"
		} else {
			if data[0] != "0" {
				res += "角"
			}
		}
	}
	return res
}

// ToUpperChinese 整数转中文大写
func ToUpperChinese(n int) string {
	var res []string
	w := (n % hundredMillion) / tenThousand
	y := n / hundredMillion
	a := strings.Split(strconv.Itoa(n%tenThousand), "")
	Reverse(a)
	var flag = false
	for i, s := range a {
		if s != "0" {
			if i != 0 {
				res = append(res, chineseUnit[6+i])
			}
			res = append(res, chineseMapping[s])
			flag = false
		} else {
			// 避免重复的 `零`
			if flag {
				continue
			}
			res = append(res, "零")
			flag = true
		}
	}
	if w > 0 {
		res = append(res, ToUpperChinese(w)+"万")
	}
	if y > 0 {
		res = append(res, ToUpperChinese(y)+"亿")
	}
	Reverse(res)
	// 尾部只有零则省略
	return strings.TrimRight(strings.Join(res, ""), "零")
}

// DecimalRetainTwo 保留两位小数
func DecimalRetainTwo(input float64) float64 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", input), 64)
	return res
}

const (
	// A 就 65
	asciiOffset = 64
	// 模 26
	mod = 26
)

func ToColumnLetter(index int) string {
	var res []byte
	for {
		res = append(res, byte(index%mod+asciiOffset))
		index /= mod
		if index < mod {
			if index != 0 {
				res = append(res, byte(index%mod+asciiOffset))
			}
			break
		}
	}

	Reverse(res)
	return string(res)
}

func ToCellIndex(x, y int) string {
	return ToColumnLetter(x) + strconv.Itoa(y)
}

var g = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}

// To16 转16进制
func To16(n int) string {
	var str string
	res := IntSplit(n, 16)
	for _, re := range res {
		str = g[re] + str
	}
	return str

}

// IntSplit 字符串转数字组
func IntSplit(n int, split int) []int {

	var res []int
	for split <= n {
		res = append(res, n%split)
		n = n / split
	}
	res = append(res, n)
	return res
}

// TimeSplit 字符串转时间
func TimeSplit(str, split, format string) (times []time.Time) {
	var list = strings.Split(str, split)
	for _, s := range list {
		parse, _ := time.Parse(format, s)
		times = append(times, parse)
	}
	return
}

func Max[T float32 | float64 | int | int64](ns ...T) T {
	var maxT T
	for i, t := range ns {
		if i == 0 {
			maxT = t
		} else {
			if t > maxT {
				maxT = t
			}
		}
	}
	return maxT
}

// TrimRightZeroRetain
//
//	最多保留两位0
//	TrimRightZeroRetain("26.0583",2)    26.0583
//	TrimRightZeroRetain("26.050000000",2)    26.05
func TrimRightZeroRetain(str string, remain int) string {
	z := strings.Split(str, ".")
	if len(z) == 0 {
		return str
	}
	temp := strings.TrimRight(z[len(z)-1], "0")
	for len(temp) < remain {
		temp += "0"
	}
	z[len(z)-1] = temp
	return strings.Join(z, ".")
}
