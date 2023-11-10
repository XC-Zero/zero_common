package excel

import (
	"github.com/xuri/excelize/v2"
	"testing"
)

func TestToColumnLetter(t *testing.T) {
	file, err := excelize.OpenFile("A:\\download\\test (2).xlsx")
	if err != nil {
		panic(err)
	}

	file = AddAdditionalPrivateSheet(file, 1, 2, 3, 4, 5, 6)
	file.Save()
	ReadAddition(file)
}
