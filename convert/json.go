package convert

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"strings"
)

// JsonOmitOther 忽略掉不要的 TODO !!!!
func JsonOmitOther(json []map[string]any, keys ...string) []map[string]any {
	return json
}

func SaveAsJson(file io.Reader, path string) error {
	_, data := ToJson(file)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	GB18030ToUTF8(bytes.NewBuffer(marshal), out)

	//_, err = io.Copy(out, )
	//if err != nil {
	//	return err
	//}
	return nil
}

func ToJson(file io.Reader, hooks ...func([]map[string]any) []map[string]any) ([]string, []map[string]any) {
	all, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	col := AnyDetectColIndex(bytes.NewBuffer(all))
	cols, data := AnyToJson(bytes.NewBuffer(all), col)

	if hooks != nil {
		for i := 0; i < len(hooks); i++ {
			if hooks[i] != nil {
				data = hooks[i](data)
			}
		}
	}
	return cols, data
}

func TxtToJson(file io.Reader, columnLine int) (column []string, res []map[string]any) {
	scanner := bufio.NewScanner(file)

	var now = 1
	var l int

	for scanner.Scan() {

		if now == columnLine {
			colText := scanner.Text()
			colText = strings.ReplaceAll(colText, "-", "_")
			column = strings.Split(colText, "\t")
			l = len(column)
		}
		if now > columnLine {
			data := strings.Split(scanner.Text(), "\t")

			var temp = make(map[string]any, len(data))
			for i, s := range data {
				if i >= l {
					continue
				}
				temp[column[i]] = s
			}
			res = append(res, temp)
		}
		now++
	}

	return column, res
}

// CSVToJson TODO 可以开协程加速
func CSVToJson(file io.Reader, columnLine int, splitSymbols ...byte) (column []string, res []map[string]any) {
	scanner := bufio.NewScanner(file)

	var now = 1
	var l int

	for scanner.Scan() {

		if now == columnLine {
			colText := scanner.Text()
			colText = strings.ReplaceAll(colText, "-", "_")
			colText = strings.ReplaceAll(colText, "\\ufeff", "")
			colText = strings.ReplaceAll(colText, "\\xA8\\xB9r", "")
			colText = strings.ReplaceAll(colText, "\\xA1\\xE3", "")
			column = CSVSplit(colText, splitSymbols...)

			l = len(column)
		}
		if now > columnLine {
			nowText := scanner.Text()
			if len(strings.TrimSpace(nowText)) == 0 {
				continue
			}
			nowText = strings.ReplaceAll(nowText, "&quot;", "%%%%%")
			nowText = strings.ReplaceAll(nowText, "\\xA8\\xB9r", "")
			nowText = strings.ReplaceAll(nowText, "\\xC3\\xEB", "")
			nowText = strings.ReplaceAll(nowText, "\\xA1\\xE3", "")
			data := CSVSplit(nowText, splitSymbols...)

			var temp = make(map[string]any, len(data))
			for i, s := range data {
				if i >= l {
					continue
				}
				s = strings.ReplaceAll(s, "%%%%%", "\"")
				temp[column[i]] = s
			}
			res = append(res, temp)
		}
		now++
	}

	return column, res
}

func ExcelToJson(file io.Reader, columnLine int, sheetName ...string) (column []string, res []map[string]any) {
	excel, err := excelize.OpenReader(file)
	defer excel.Close()
	if err != nil {
		log.Printf("[ERROR] Excel to json open file failed! Error stack is %v", errors.WithStack(err))
		return nil, nil
	}
	var sheet string
	if len(sheetName) != 0 {
		sheet = sheetName[0]
	} else {
		sheet = excel.GetSheetMap()[excel.GetActiveSheetIndex()]
		if sheet == "" {
			sheet = "Sheet1"
		}
	}
	rows, err := excel.GetRows(sheet)
	if err != nil {
		log.Printf("[ERROR] Excel to json failed! Error stack is %v", errors.WithStack(err))
		return nil, nil
	}
	for i, row := range rows {
		if i < columnLine {
			continue
		}
		if i == columnLine {
			column = make([]string, 0, len(row))
			for _, cell := range row {
				column = append(column, strings.ReplaceAll(cell, "-", "_"))
			}
			continue
		}
		var inner = make(map[string]any, len(column))
		for j, s := range row {
			if j < len(column) {
				inner[column[j]] = s
			}
		}

		res = append(res, inner)

	}
	return
}

// AnyToJson 任意格式文件转json格式 TODO JSON!!!!
func AnyToJson(file io.Reader, columnLine int) (column []string, res []map[string]any) {
	all, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	fileType := DetectFileType(bytes.NewBuffer(all))
	switch {
	case IsExcel(fileType):
		return ExcelToJson(bytes.NewBuffer(all), columnLine)
	case fileType == TXT:
		return TxtToJson(bytes.NewBuffer(all), columnLine)
	case fileType == CSV:
		return CSVToJson(bytes.NewBuffer(all), columnLine)
	case fileType == JSON:

	}
	return TxtToJson(bytes.NewBuffer(all), columnLine)
}
