package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"time"
)

const tessanSheet = "tessan_private"

func CopyExcel() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

}

// AddAdditionalPrivateSheet 添加隐藏Sheet  FIXME 老是添加失败
func AddAdditionalPrivateSheet[T any](excel *excelize.File, data ...T) *excelize.File {
	hiddenStyle, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "FFFFFF",
		},
		Protection: &excelize.Protection{Locked: true},
	})
	if err != nil {
		log.Println(err)
	}
	index := -1
	for i, name := range excel.GetSheetMap() {
		if name == tessanSheet {
			index = i
		}
	}
	if index == -1 {
		index, err = excel.NewSheet(tessanSheet)
		if err != nil {
			log.Println(err)
		}
	}
	excel.SetCellValue(tessanSheet, "H1", "data")
	excel.SetCellStyle(tessanSheet, "H1", "H1", hiddenStyle)

	for i, datum := range data {
		cellIndex, _ := excelize.CoordinatesToCellName(8, i+2)
		err = excel.SetCellValue(tessanSheet, cellIndex, datum)
		if err != nil {
			log.Println(err)
		}
		err = excel.SetCellStyle(tessanSheet, cellIndex, cellIndex, hiddenStyle)
		if err != nil {
			log.Println(err)
		}
	}

	excel.SetCellValue(tessanSheet, "C1", "generate_time")
	excel.SetCellValue(tessanSheet, "C2", time.Now().Format("2006/01/02 15:04 05"))
	err = excel.SetSheetVisible(tessanSheet, false)
	if err != nil {
		log.Println(err)
	}
	err = excel.Save()
	if err != nil {
		log.Println(err)
	}
	return excel
}

func ReadAddition(excel *excelize.File) []any {
	for i, s := range excel.GetSheetMap() {
		log.Println(i, s)
	}
	return []any{}
}
