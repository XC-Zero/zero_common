package convert

import (
	"bytes"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"testing"
)

func TestAnyDetectColIndex(t *testing.T) {
	open, err := os.Open("A:\\xccjm\\Documents\\a.txt")
	if err != nil {
		panic(err)
	}
	create, err := os.Create("A:\\xccjm\\Documents\\u.txt")
	if err != nil {
		return
	}
	GB18030ToUTF8(open, create)

}

func TestGB18030ToUTF8(t *testing.T) {

	res := bytes.NewBufferString("")
	Transform(bytes.NewBuffer([]byte("ó")), GB18030, res, UTF8)
	a := res.String()
	log.Println(a)

	r := bytes.NewBufferString("")
	Transform(bytes.NewBuffer([]byte(a)), UTF8, r, GB18030)
	log.Println(r.String())

}

//
//func TestGB2312ToUTF8(t *testing.T) {
//	all, _ := io.ReadAll(GB2312ToUTF8(bytes.NewBuffer([]byte("車"))))
//	log.Println(string(all))
//	all, _ = io.ReadAll(GB2312ToUTF8(bytes.NewBuffer([]byte("ó"))))
//	log.Println(string(all))
//}

func TestSliceConvert(t *testing.T) {
	file, err := os.Open("A:\\xccjm\\Documents\\test01x.xlsx")
	if err != nil {
		panic(err)
	}
	err = SaveAsJson(file, "A:\\xccjm\\Documents\\test01x.json")
	if err != nil {
		panic(err)
	}
}

func TestCSVDetectColumnIndex(t *testing.T) {
	open, err := os.Open("D:\\tessan_erp技术文档\\比对报告\\日期范围报告\\Tessan 2023-5月\\ca\\Tessan ca20235月MonthlyTransaction.csv")
	if err != nil {
		panic(err)
	}
	a := AnyDetectColIndex(open)
	log.Println(a)
}

func TestCSVUtil(t *testing.T) {
	open, err := os.Open("D:\\tessan_erp技术文档\\比对报告\\日期范围报告\\Tessan 2023-5月\\us\\Tessan us20235月MonthlyUnifiedTransaction.csv")
	if err != nil {
		panic(err)
	}

	var m []map[string]any
	if err := gocsv.UnmarshalFile(open, &m); err != nil { // Load clients from file
		panic(err)
	}
	//spew.Dump(m[:10])
}

func TestCSVSplit(t *testing.T) {
	//tt := CSVSplit(`"All amounts in USD, unless specified",,,,,,,,,,,,,,,,,,,,,,,,,,,,,`)
	//t2 := CSVSplit(`date/time,settlement id,type,order id,sku,description,quantity,marketplace,account type,fulfillment,order city,order state,order postal,tax collection model,product sales,product sales tax,shipping credits,shipping credits tax,gift wrap credits,giftwrap credits tax,Regulatory Fee,Tax On Regulatory Fee,promotional rebates,promotional rebates tax,marketplace withheld tax,selling fees,fba fees,other transaction fees,other,total`)
	//t3 := CSVSplit(`"May 1, 2023 12:00:16 AM PDT",17833217601,Order,112-5296892-7513800,TS-165-C-GR-US,"Multi Plug Outlet Extender, USB Wall Charger, TESSAN Multiple Outlet Splitter with 5 Outlets and 3 USB (1 USB C), Electrical Power Expander with Surge",1,amazon.com,Standard Orders,Amazon,SAN DIEGO,CA,92115-6801,MarketplaceFacilitator,25.99,2.01,2.99,0,0,0,0,0,-2.99,0,-2.01,-3.9,-3.86,0,0,18.23,`)
	//

}

func TestDetectFileType(t *testing.T) {
	open, err := os.Open("D:\\tessan_erp技术文档\\比对报告\\日期范围报告\\Tessan 2023-5月\\us\\Tessan us20235月MonthlyUnifiedTransaction.csv")
	if err != nil {
		panic(err)
	}

	log.Println(DetectFileType(open))
}
