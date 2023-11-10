package convert

import (
	"github.com/davecgh/go-spew/spew"
	"os"
	"testing"
)

func TestToJson(t *testing.T) {
	open, err := os.Open("C:\\Users\\xccjm\\Desktop\\虾皮-富创兴 MX 订单.xlsx")
	if err != nil {
		panic(err)
	}
	spew.Dump(ToJson(open))
}

func TestExcelToJson(t *testing.T) {
	open, err := os.Open("C:\\Users\\xccjm\\Desktop\\虾皮-富创兴 TW1订单.xlsx")
	if err != nil {
		panic(err)
	}
	spew.Dump(ExcelToJson(open, 1, "orders"))
}
