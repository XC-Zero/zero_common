package convert

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"regexp"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
	"io"
	"log"
	"strings"
	"time"
)

func ToParquet(file io.Reader) {

}

func JsonToParquet(data []map[string]any) {

}

// TxtToParquet FIXME 目前仅取第一行的数据作为类型
func TxtToParquet(file io.Reader, fileName string, columnLine int) (column []string) {
	scanner := bufio.NewScanner(file)
	var now = 1
	var l int
	var err error

	fw, err := local.NewLocalFileWriter((fileName) + ".parquet")
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	//write
	var pw *writer.CSVWriter
	var baseMd []string
	var colType = make(map[string]reflect.Type)
	for scanner.Scan() {

		if now == columnLine {
			colText := scanner.Text()
			colText = strings.ReplaceAll(colText, "-", "_")
			temp := strings.Split(colText, "\t")
			for _, s := range temp {
				now := strings.Trim(s, " ")
				if len(now) == 0 {
					continue
				}
				column = append(column, now)

			}
			l = len(column)
			log.Println("col is :", column, l)

		}

		if now == columnLine+1 {
			t := splitText(scanner.Text())

			for i, s := range column {
				aha, tt := getParquetKeyType(t[i])
				baseMd = append(baseMd, fmt.Sprintf(`name=%s,type=%s`, s, aha))
				colType[s] = tt
			}
			log.Println(baseMd)

			pw, err = writer.NewCSVWriter(baseMd, fw, 4)
			if err != nil {
				log.Printf("Can't create parquet writer :%+v", errors.WithStack(err))

			}
		}
		if now > columnLine {
			t := splitText(scanner.Text())

			data := make([]any, 0, len(t))
			for i, s := range t {
				t, ok := colType[column[i]]
				if !ok {
					continue
				}
				data = append(data, anyToType(stringToAny(s), t))
			}

			if now%10000 == 2 {
				log.Println("data is :", data, len(data))
			}

			err := pw.Write(data)
			if err != nil {
				log.Println(data)
				panic(err)
			}
		}
		now++
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
		return
	}
	fw.Close()

	return
}

func anyToType(data any, t reflect.Type) any {
	v, nt := reflect.ValueOf(data), reflect.TypeOf(data)
	if nt == t {
		return data
	}
	if v.CanConvert(t) {
		return v.Convert(t).Interface()
	}
	if t == reflect.TypeOf("") {
		return fmt.Sprintf("%v", data)
	}
	return reflect.Zero(t)
}

func getParquetKeyType(data any) (string, reflect.Type) {
	switch data.(type) {
	case string:
		return "BYTE_ARRAY", reflect.TypeOf("")
	case float64:
		return "DOUBLE", reflect.TypeOf(float64(0))
	case float32:
		return "FLOAT", reflect.TypeOf(float32(0))
	case int64, int:
		return "INT64", reflect.TypeOf(int64(0))
	case int32, int8, int16:
		return "INT32", reflect.TypeOf(int32(0))
	case bool:
		return "BOOLEAN", reflect.TypeOf(true)
	default:
		return "BYTE_ARRAY", reflect.TypeOf("")
	}
}

func getParquetType[T any | string | int | float64](data any) string {
	switch data.(type) {
	case time.Time:
		return "BYTE_ARRAY"
	case string:
		return "FIXED_LEN_BYTE_ARRAY"
	case float64:
		return "DOUBLE"
	case float32:
		return "FLOAT"
	case int64, int:
		return "INT64"
	case int32, int8, int16:
		return "INT32"
	case bool:
		return "BOOLEAN"
	case map[string]T, []T:
		return "MAP"
	default:
		return "BYTE_ARRAY"
	}
}

func getFmtType(data any) string {
	switch data.(type) {
	case string:
		return `"%s"`
	case float64:
		return "%.2f"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "%d"
	case bool:
		return "%t"
	default:
		return "%v"
	}
}

func stringToAny(datum string) any {
	if ok, _ := regexp.MatchString(timeMatch, datum); ok {
		if parse, err := time.Parse(time.RFC3339, datum); err != nil {
			return parse.Format("2006-01-02 15:04:05")
		}
	}
	if strings.ToLower(datum) == "false" {
		return false
	}
	if strings.ToLower(datum) == "true" {
		return true
	}

	if IsVailAmount(datum) {
		return NewAmount(datum).ParseFloat64()
	}
	return datum
}

const timeMatch = `^((\d){4}(-|_|\\|/| )(\d){1,2}(-|_|\\|/| )(\d){1,2})(.*)`

func stringFormat(datum string) string {

	if ok, _ := regexp.MatchString(timeMatch, datum); ok {
		if parse, err := time.Parse(time.RFC3339, datum); err != nil {
			return parse.Format("2006-01-02 15:04:05")
		}
	}

	if strings.ToLower(datum) == "false" {
		return "false"
	}
	if strings.ToLower(datum) == "true" {
		return "true"
	}

	if IsVailAmount(datum) {
		return fmt.Sprintf("%.2f", NewAmount(datum).ParseFloat64())
	}
	return datum
}
