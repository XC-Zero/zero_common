package main

//
//import (
//	"gitlab.tessan.com/data-center/tessan-erp-model/mysql"
//	"gitlab.tessan.com/data-center/tessan-erp-model/ts_mid/commodity"
//	"gitlab.tessan.com/data-center/tessan-erp-model/ts_mid/elementary_info"
//)
//
//var anys = []any{
//	&elementary_info.Msku{},
//	&elementary_info.Sku{},
//	&elementary_info.Provider{},
//	&elementary_info.Store{},
//	&elementary_info.Platform{},
//	&elementary_info.Location{},
//	&elementary_info.Subsidiary{},
//	&elementary_info.Brand{},
//	&elementary_info.SkuProviderRelation{},
//	&elementary_info.SkuMskuRelation{},
//	&elementary_info.SkuLocationRelation{},
//
//	&commodity.ReturnOrder{},
//	&commodity.BatchPlan{},
//	&commodity.TransferOrder{},
//	&commodity.CustomsDeclaration{},
//	&commodity.TransferOrder{},
//	&commodity.TransferOrderDetail{},
//	&commodity.ReceiptOrder{},
//	&commodity.ReceiptOrderDetail{},
//	&commodity.PurchaseOrder{},
//	&commodity.PurchaseOrderDetail{},
//}
//
//func main() {
//	err := mysql.GetLocalDBInstance(true).AutoMigrate(anys...)
//	if err != nil {
//		panic(err)
//	}
//}
