package excel

import (
	"github.com/xuri/excelize/v2"
	"reflect"
)

var (
	defaultFont = excelize.Font{
		Bold:         false,
		Italic:       false,
		Underline:    "",
		Family:       "等线",
		Size:         14,
		Strike:       false,
		Color:        "",
		ColorIndexed: 0,
		ColorTheme:   nil,
		ColorTint:    0,
		VertAlign:    "",
	}
	defaultBorder = []excelize.Border{
		{
			Type:  "left",
			Color: "#000000",
			Style: 1,
		},
		{
			Type:  "right",
			Color: "#000000",
			Style: 1,
		},
		{
			Type:  "top",
			Color: "#000000",
			Style: 1,
		},
		{
			Type:  "bottom",
			Color: "#000000",
			Style: 1,
		},
	}

	// EmptyStyle 字段为空时  样式
	EmptyStyle = excelize.Style{
		Border: []excelize.Border{{
			Type:  "bottom",
			Color: "#FD5151",
			Style: 5,
		}},
	}

	DefaultStyle = excelize.Style{}

	// FillStripeBlue1 斑马纹样式 1
	FillStripeBlue1 = excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Shading: 5,
			Color:   []string{"#DDEBF7"},
		}, Font: &defaultFont, Border: defaultBorder}

	// FillStripeBlue2 斑马纹样式 2
	FillStripeBlue2 = excelize.Style{Fill: excelize.Fill{
		Type:    "pattern",
		Pattern: 1,
		Shading: 5,
		Color:   []string{"#DAE3F3"},
	}, Font: &defaultFont, Border: defaultBorder}

	OneSizeFont = excelize.Style{Font: &excelize.Font{Size: 1, Color: "#FFFFFF"}}
)

// SetValueWithStyle 设置样式和值
func SetValueWithStyle(file *excelize.File, sheet, cell string, value any, styleID int, onlyEmpty bool) {

	file.SetCellValue(sheet, cell, value)

	if onlyEmpty {
		if reflect.ValueOf(value).IsZero() {
			file.SetCellStyle(sheet, cell, cell, styleID)
		}
		return
	}
	file.SetCellStyle(sheet, cell, cell, styleID)
	return
}
