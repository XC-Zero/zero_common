package convert

import (
	"bufio"
	"bytes"
	"github.com/goccy/go-json"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultSheet = "Sheet1"
const maxDetectCol = 50

type FileType string

const (
	Unknown FileType = ""
	CSV     FileType = "data:text/csv"
	TXT     FileType = "application/txt"
	TXT2    FileType = "text/plain"
	JSON    FileType = "application/json"
	XML     FileType = "application/xml"
	ZIP     FileType = "application/zip"
	GZIP    FileType = "application/x-gzip"
	XLSX    FileType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,"
	XLS     FileType = "application/vnd.ms-excel"
	XLSM    FileType = "application/vnd.ms-excel.sheet.macroEnabled.12"
	OCT     FileType = "application/octet-stream"
	ICO     FileType = "image/vnd.microsoft.icon"
	WASM    FileType = "application/wasm"
	PNG     FileType = "image/png"
)

func DetectFileType(file io.Reader) FileType {
	all, err := io.ReadAll(file)
	if err != nil {
		return Unknown
	}

	var fileType = temp(http.DetectContentType(all))
	switch {
	case fileType.Contains(JSON):
		return JSON
	case fileType.Contains(TXT), fileType.Contains(TXT2):
		if json.Valid(all) {
			return JSON
		}
		var str string
		if len(all) < 1000 {
			str = string(all)
		} else {
			str = string(all[:1000])
		}
		strLen := len(str)
		// 如果逗号数量是tab数量的一倍还多,认为是csv FIXME 其实还应该看看每行的分隔符是不是一样多的
		point := strLen - len(strings.ReplaceAll(str, ",", ""))
		tab := strLen - len(strings.ReplaceAll(str, "\t", ""))
		if point/2 > tab {
			return CSV
		}
		return TXT

	case fileType.Contains(CSV):
		return CSV
	case fileType.Contains(ZIP):
		return ZIP
	case fileType.Contains(XLS):
		return XLS
	case fileType.Contains(XLSM):
		return XLSM
	case fileType.Contains(XLSX):
		return XLSX
	case fileType.Contains(OCT):
		return OCT
	case fileType.Contains(PNG):
		return PNG
	case fileType.Contains(GZIP):
		return GZIP
	case fileType.Contains(XML):
		return XML
	case fileType.Contains(WASM):
		return WASM
	default:
		return Unknown
	}
}

// TxtToCSV 文本转csv
func TxtToCSV(file io.Reader) *excelize.File {
	excel := excelize.NewFile()
	excel.NewSheet(defaultSheet)
	scanner := bufio.NewScanner(file)
	var row = 1
	for scanner.Scan() {
		column := strings.Split(scanner.Text(), "\t")
		for i, s := range column {
			idx, _ := excelize.CoordinatesToCellName(i+1, row)
			excel.SetCellStr(defaultSheet, idx, s)
		}
		row++
	}
	return excel
}

func FirstLine(file io.Reader) (res string) {
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		res = scanner.Text()

	}
	return
}

func NthLine(file io.Reader, n int) (res string) {
	scanner := bufio.NewScanner(file)
	v := 0
	for scanner.Scan() {
		v++
		if v == n {
			return scanner.Text()
		}
	}
	return
}

// TxtDetectColumnIndex 探测col位置
// 根据tab符剪切 最长的第一行就是头
func TxtDetectColumnIndex(file io.Reader) int {

	scanner := bufio.NewScanner(file)
	log.Println("Success read txt file!")

	var maxLen, firstIndex int
	for i := 0; i < maxDetectCol; i++ {
		if scanner.Scan() {
			l := len(strings.Split(scanner.Text(), "\t"))
			if l > maxLen {
				maxLen = l
				firstIndex = i
			}
		}
	}
	return firstIndex + 1
}

// CSVDetectColumnIndex 探测col位置
// 根据 , 符剪切 最长且非空最多的第一行就是头
func CSVDetectColumnIndex(file io.Reader) int {

	scanner := bufio.NewScanner(file)

	var maxLen, firstIndex int
	for i := 0; i < maxDetectCol; i++ {
		if scanner.Scan() {
			data := CSVSplit(scanner.Text())
			var notN = 0
			for _, s := range data {
				if s != "" {
					notN++
				}
			}
			if notN > maxLen {
				maxLen = notN
				firstIndex = i
			}
		}
	}
	return firstIndex + 1
}

func ExcelDetectColumnIndex(file io.Reader) int {
	excel, err := excelize.OpenReader(file)
	if err != nil {
		log.Println(err)
		return 0
	}
	var maxLen, firstIndex int
	var sheet = excel.GetSheetMap()[0]

	rows, err := excel.GetRows(sheet)
	if err != nil {
		log.Println(err)
		return 0
	}

	for i := 0; i < maxDetectCol; i++ {
		l := len(rows[i])
		if l > maxLen {
			firstIndex = i
		}
	}

	return firstIndex + 1
}

func AnyDetectColIndex(file io.Reader) int {
	all, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return 0
	}
	fileType := DetectFileType(bytes.NewBuffer(all))
	switch {
	case IsExcel(fileType):
		return ExcelDetectColumnIndex(bytes.NewBuffer(all))
	case fileType == CSV:
		return CSVDetectColumnIndex(bytes.NewBuffer(all))
	case fileType == TXT:
		return TxtDetectColumnIndex(bytes.NewBuffer(all))

	}

	return 0
}

type temp string

func (h temp) Contains(fileType FileType) bool {
	return strings.Contains(string(h), string(fileType))
}

func IsExcel(fileType FileType) bool {
	switch fileType {
	case ZIP, XLSM, XLS, XLSX, OCT:
		return true
	default:
		return false
	}
}

func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// mkdir 创建目录，mkdirAll 可创建多层级目录
		os.MkdirAll(path, os.ModePerm)
	}
}

//// CSVSplit 利用带下标的符号栈,分割
//func CSVSplit(row string) []string {
//	var flagIdx = make(map[byte]int)
//	var res []string
//	var firstEmpty = true
//	var i = 0
//	for {
//		if row[i] == ',' {
//			if i == 0 {
//				res = append(res, "")
//				goto LoopPoint
//			}
//			//如果这个逗号之前出现了一个双引号, 逗号栈清空
//			if p, ok := flagIdx['"']; ok && p < i {
//				delete(flagIdx, ',')
//				goto LoopPoint
//
//			} else if !ok {
//				// 如果这个逗号之前还有一个逗号
//				if lastIdx, okk := flagIdx[',']; okk {
//					if lastIdx == i-1 {
//						res = append(res, "")
//					} else {
//						res = append(res, row[lastIdx+1:i])
//					}
//					flagIdx[','] = i
//					goto LoopPoint
//					// 如果前面啥也没有,且第一次
//				} else if firstEmpty {
//					res = append(res, row[:i])
//					flagIdx[','] = i
//					firstEmpty = false
//					goto LoopPoint
//				} else {
//					flagIdx[','] = i
//					goto LoopPoint
//				}
//			}
//
//		}
//		if row[i] == '"' {
//			firstEmpty = false
//			if lastIdx, ok := flagIdx['"']; ok {
//				if lastIdx == i-1 {
//					res = append(res, "")
//				} else {
//					res = append(res, row[lastIdx+1:i])
//				}
//				// 出栈
//				delete(flagIdx, '"')
//
//			} else {
//				// 入栈
//				flagIdx['"'] = i
//			}
//			// 如果下一个是逗号 ,清空逗号栈
//			if i+1 < len(row) && row[i+1] == ',' {
//				delete(flagIdx, ',')
//			}
//		}
//	LoopPoint:
//		i++
//		if i == len(row) {
//			if lastIdx, okk := flagIdx[',']; okk {
//				d := row[lastIdx+1:]
//				if !strings.Contains(d, "\"") {
//					res = append(res, d)
//				}
//
//				delete(flagIdx, ',')
//			}
//			return res
//		}
//	}
//}

type symbol struct {
	symbol byte
	idx    int
}

type symbolStack struct {
	stack   []symbol
	lastIdx map[byte]int
}

func (s *symbolStack) push(sym symbol) int {
	if s.lastIdx == nil {
		s.lastIdx = make(map[byte]int)
	}

	// 压栈
	s.stack = append(s.stack, sym)

	if idx, ok := s.lastIdx[sym.symbol]; ok {
		// 出栈
		for i, s2 := range s.stack {
			if s2.idx == idx {
				if i == 0 {
					s.stack = []symbol{}
				} else {
					s.stack = s.stack[:i-1]
				}
			}
		}
		delete(s.lastIdx, sym.symbol)
		return idx
	} else {
		s.lastIdx[sym.symbol] = sym.idx
		return -1
	}
}

// CSVSplit 利用带下标的符号栈,分割
func CSVSplit(row string) []string {
	row = "," + row + ","
	var symbolStack symbolStack
	var res []string
	var idxs []int
	var i = 0
	for {
		sym := symbol{
			symbol: row[i],
			idx:    i,
		}
		if row[i] == ',' {
			if _, ok := symbolStack.lastIdx['"']; ok {
				goto MainLoop
			}
			if idx := symbolStack.push(sym); idx != -1 {
				idxs = append(idxs, idx+1, i)
			}
		} else if row[i] == '"' {
			symbolStack.push(sym)
		}

	MainLoop:
		i++
		if i == len(row) {
			idxs = append(idxs, i-1)
			for j := 1; j < len(idxs); j++ {
				res = append(res, strings.Trim(row[idxs[j-1]:idxs[j]], `", `))
			}
			break
		}
	}
	return res
}

func NewTimeNameFile(data []byte) {
	file, err := os.Create(time.Now().Format("2006-01-02 15:04 05"))
	if err != nil {
		return
	}

	io.Copy(file, bytes.NewReader(data))
	file.Close()
}
