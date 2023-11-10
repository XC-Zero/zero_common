package convert

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// MaybeTime 探测是否可能为时间格式
func MaybeTime(str string) bool {
	str = strings.ReplaceAll(str, "-", " ")
	str = strings.ReplaceAll(str, "−", " ")
	str = strings.ReplaceAll(str, "/", " ")
	str = strings.ReplaceAll(str, "\\", " ")
	str = strings.ReplaceAll(str, " ", " ")
	str = strings.ReplaceAll(str, " ", " ")
	str = strings.ReplaceAll(str, " ", " ")
	str = strings.ReplaceAll(str, ",", " ")
	str = strings.ReplaceAll(str, ".", " ")
	str = strings.ReplaceAll(str, "   ", " ")
	str = strings.ReplaceAll(str, "  ", " ")
	block := strings.Split(str, " ")
	// 起码要有 年 月 日 时分秒  四个块
	// 目前所见最长的就  年 月 日 时 分 秒 上下午 时区
	blockLen := len(block)
	if blockLen < 4 || blockLen > 8 {
		return false
	}
	var colon, number, strN, avgLen int
	avgLen = (len(str) / blockLen) - 1
	// 平均每块都很长
	if avgLen > 8 {
		return false
	}
	for _, s := range block {
		if strings.Contains(s, ":") {
			colon++
			continue
		}
		if _, err := strconv.Atoi(s); err == nil {
			number++
			continue
		}
		strN++
	}

	// 全是数字
	if number == blockLen {
		return false
	}
	// 全是字符串
	if strN == blockLen {
		return false
	}
	// 俩有冒号的块
	if colon > 1 {
		return false
	}

	return true
}

type TimeLayout string

var (
	month = map[string][]string{
		"Jan": {"Jan", "janv", "Janvier", "Ene", "Enero", "gen", "gennaio", "jan", "januari", "sty", "styczeń", "januari", "Ocak", "Oca"},
		"Feb": {"Feb", "Febrero", "Fév", "févr", "Févier", "febbraio", "februari", "lut", "luty", "Şubat", "Şub"},
		"Mar": {"Mar", "Mars", "Mart", "mar", "mars", "marzec", "mrt", "mrt.", "maart", "Marzo", "Mzo"},
		"Apr": {"Apr", "Avr", "Avril", "abr", "avr", "apr", "apr.", "avr.", "Abr", "Abril", "aprile", "april", "kwi", "kwiecień", "Nisan", "Nis"},
		"May": {"May", "Mayıs", "ma", "maj", "me", "mei", "mag", "maggio", "My", "Mayo", "Mai", "mai", "mai."},
		"Jun": {"Jun", "Juin", "jun.", "juin", "juin.", "Junio", "giu", "giugno", "juni", "cze", "czerwiec", "Haziran", "Haz"},
		"Jul": {"Jul", "Julio", "Juil", "Juillet", "Temmuz", "Tem", "jul", "jul.", "juil.", "juli", "lip", "lipiec", "lug", "luglio"},
		"Aug": {"Aug", "Août", "août", "August", "Ago", "Agosto", "ago", "aug", "aug.", "augustus", "sie", "sierpień", "augusti", "Ağustos", "Ağu"},
		"Sep": {"Sep", "Sept", "sept.", "Septembre", "Septiembre", "set", "sep", "sep.", "settembre", "sept", "september", "wrz", "wrzesień", "Eylül", "Eyl"},
		"Oct": {"Oct", "oct", "oct.", "okt.", "Ekim", "Eki", "okt", "oktober", "paź", "październik", "ott", "ottobre", "Octubre", "Octobre"},
		"Nev": {"Nov", "nov", "nov.", "November", "Novembre", "Noviembre", "novembre", "november", "lis", "listopad", "november", "Kasım", "Kas"},
		"Dec": {"Dec", "Aralık", "Ara", "dec", "dec.", "december", "gru", "grudzień", "Dic", "dicembre", "Diciembre", "Déc", "Décembre", "déc"},
	}
	MonthMapping map[string]string
	AMMapping    = map[string]string{
		"p.m.": "PM",
		"a.m.": "AM",
		"am":   "AM",
		"pm":   "PM",
		"p.m":  "PM",
		"a.m":  "AM",
		"Am":   "AM",
		"Pm":   "PM",
	}
)

func init() {
	MonthMapping = make(map[string]string)
	for k, s := range month {
		for _, asia := range s {
			MonthMapping[asia] = k
		}
	}
}

func TransformTime(layout, str string) (time.Time, error) {
	str = strings.ReplaceAll(str, "−", "-")
	str = strings.ReplaceAll(str, " ", " ")
	str = strings.ReplaceAll(str, " ", " ")
	var b = strings.Split(str, " ")
	for i, s := range b {
		if v, ok := MonthMapping[s]; ok {
			b[i] = v
		}
		if v, ok := AMMapping[s]; ok {
			b[i] = v
		}
	}

	str = strings.Join(b, " ")
	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func TryToLocalTime(str, layout string) string {

	local, err := TransformTime(str, layout)
	if err != nil {
		return str
	}
	return local.Format(time.RFC3339)
}

func ParseStr(str string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		log.Println("time parse Error:", err)
	}
	return t
}

// AbandonTimeZone 舍弃时区的时间 (改时区不改时间)
func AbandonTimeZone(date time.Time) time.Time {
	_, offset := date.Zone()
	date = date.Add(time.Duration(offset) * time.Second)
	return date.UTC()
}

// ReverseTime 取相反时区时间!
func ReverseTime(date time.Time, location *time.Location) time.Time {
	t := date.In(location)
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
}
