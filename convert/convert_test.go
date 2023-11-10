package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"unsafe"
)

func TestGetAnyStructTag(t *testing.T) {
	var a = []struct {
		Hello string `label:"你好"`
	}{{Hello: "你好"}}
	var b []int
	log.Println(a)
	log.Println(GetAnyStructTag(a, "label"))
	log.Println(GetAnyStructTag(b, "label"))

}

func TestTxtToCSV(t *testing.T) {

	open, err := os.Open("A:\\xccjm\\Documents\\b.txt")
	if err != nil {
		panic(err)
	}
	err = TxtToCSV(open).SaveAs("A:\\xccjm\\Documents\\b.xlsx")
	if err != nil {
		panic(err)
	}
}

func TestTxtToJson(t *testing.T) {

	open, err := os.Open("A:\\xccjm\\Documents\\(2023-07-05 09-36-49)TESSAN 美国-6月.txt")
	if err != nil {
		panic(err)
	}
	column, data := TxtToJson(open, 0)

	log.Println(column)
	log.Println(unsafe.Sizeof(data), unsafe.Sizeof(data[0]))
	log.Println(unsafe.Sizeof(""))
	log.Println(unsafe.Sizeof(map[string]string{}))
}

func TestToPrimaryMap(t *testing.T) {
	open, err := os.Open("A:\\xccjm\\Documents\\a.txt")
	if err != nil {
		panic(err)
	}
	column, data := TxtToJson(open, 0)
	var i = 0
	for s, m := range ToPrimaryMap(data, column[0]) {
		spew.Dump(s, m)
		if i > 3 {
			break
		}
		i++
	}
}

func TestDifferent(t *testing.T) {

	//client, err := mysql.InitClient(config.MysqlConfig{
	//	Host:     "NAS",
	//	Port:     "49197",
	//	Username: "root",
	//	Pass:     "MySQL.123",
	//	DBName:   "amazon",
	//})
	//if err != nil {
	//	panic(err)
	//}
	wg := sync.WaitGroup{}
	wg.Add(2)

	var pa, pb PrimaryMap

	go func() {

		open2, err := os.Open("A:\\xccjm\\Documents\\d.txt")
		if err != nil {
			panic(err)
		}
		_, b := TxtToJson(open2, 1)
		//b:=getData(client)
		pb = ToPrimaryMap(b, "order_id")
		wg.Done()
		log.Println("db read done!")

	}()

	go func() {
		open2, err := os.Open("A:\\xccjm\\Documents\\a.txt")
		if err != nil {
			panic(err)
		}
		_, a := TxtToJson(open2, 1)
		pa = ToPrimaryMap(a, "amazon-order-id")
		wg.Done()
		log.Println("file parse done!")

	}()
	wg.Wait()
	dd := Different(pa, pb, map[string]string{
		"amazon-order-id": "order_id",
		"sku":             "msku",
	})
	log.Println("dd", len(dd))
	create, err := os.Create("./temp.json")
	if err != nil {
		panic(err)
	}

	marshal, err := json.Marshal(dd)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	json.Indent(&out, marshal, "", "    ")

	write, err := out.WriteTo(create)
	if err != nil {
		panic(err)
	}
	log.Println(write)

}

func getData(client *gorm.DB) []map[string]any {
	var count int64
	var wg sync.WaitGroup
	baseSql := `
select concat_ws('-', LPAD(order_id_1, 3, '0'), LPAD(order_id_2, 7, '0'), LPAD(order_id_3, 7, '0')) as order_id,
       id,
       msku_id,
       payment_datetime,
       settlement_id,
       fulfillment_id,
       tax_model_id,
       order_datetime,
       quantity,
       product_sales,
       product_sales_tax,
       shipping_credits,
       shipping_credits_tax,
       giftwrap_credits,
       giftwrap_credits_tax,
       promotional_rebates,
       promotional_rebates_tax,
       marketplace_withheld_tax,
       selling_fees,
       fba_fees,
       other_transaction_fees,
       other,
       total,
       sales_tax_collected,
       order_state,
       add_vat,
       difference,
       insert_date
from amazon.tb_payment_order
where order_datetime >= '2023-06-01' and  order_datetime < '2023-06-10'
`
	err := client.Table("tb_payment_order").
		Where(`order_datetime >= '2023-06-01' and  order_datetime < '2023-06-10'`).
		Count(&count).Error
	if err != nil {
		panic(err)
	}
	log.Println(count)
	var res = make([]map[string]any, 0, count)
	var lock sync.Mutex
	var batchNum = 100 * 100
	var times = int(count) / batchNum
	if int(count)%batchNum > 0 {
		times++
	}
	wg.Add(times)
	for i := 0; i < times; i++ {
		go func(i int) {
			var temp = make([]map[string]any, 0, batchNum)
			start, end := (i)*batchNum, (i+1)*batchNum
			sql := baseSql + fmt.Sprintf(` limit %d,%d`, start, end)
			err2 := client.Raw(sql).Scan(&temp).Error
			if err2 != nil {
				panic(err)
			}
			lock.Lock()
			res = append(res, temp...)
			lock.Unlock()
			wg.Done()
			log.Printf("times %d is done!", i)

		}(i)
	}
	wg.Wait()
	return res
}

func TestXxx(t *testing.T) {
	var prefix = "C:\\Users\\xccjm\\Desktop\\"
	tttypes(prefix + "xlsm.xlsm")
	tttypes(prefix + "xlsx.xlsx")
	tttypes(prefix + "xls.xls")
	tttypes(prefix + "csv.csv")
	tttypes(prefix + "txt.txt")
	tttypes(prefix + "xx.json")
	tttypes(prefix + "xxx.yml")

}
func tttypes(path string) {
	open, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	all, err := io.ReadAll(open)
	if err != nil {
		panic(err)
	}
	log.Println(http.DetectContentType(all))
}

func TestTransformToUTF8(t *testing.T) {
	//gb := mahonia.NewDecoder("gb18030")
	s, err := simplifiedchinese.GB18030.NewEncoder().String("\u0026")
	if err != nil {
		panic(err)
	}
	log.Println([]byte(s))
	all, err := simplifiedchinese.GB18030.NewDecoder().String(s)
	if err != nil {
		panic(err)
	}
	log.Println(all)
	//log.Println(TransformToUTF8(bytes.NewBufferString("好美丽查询")))
}

func TestSortBy(t *testing.T) {
	//var slice []int
	//var a1 = make([]int, 0)
	//var b1 = make([]int, 0)
	//log.Printf("%+v \n %+v \n %+v \n",
	//	*(*reflect.SliceHeader)(unsafe.Pointer(&slice)),
	//	*(*reflect.SliceHeader)(unsafe.Pointer(&a1)),
	//	*(*reflect.SliceHeader)(unsafe.Pointer(&b1)))
	//log.Println(&a1 == &b1, &a1, &b1)
	//log.Println(&slice == &b1)
	var a, b = struct{}{}, struct{}{}
	log.Println(&a == &b)

	//log.Printf("a %p b %p", &a, &b)

}

func TestXXX(t *testing.T) {
	a := []int{0, 1, 2}
	var b = make(map[int]*int)
	for i, v := range a {
		b[i] = &v
	}
	for key, val := range b {
		log.Println(key, "->", *val)
	}

}
