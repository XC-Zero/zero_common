package main

//import (
//	"context"
//	"github.com/davecgh/go-spew/spew"
//	"github.com/tessan/tessan_erp_proto/generate/server/goods"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//)
//
////go:generate go build -o tx-erp ./main.go
//
//func main() {
//	client, err := grpc.DialContext(context.Background(), "localhost:8087", grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		panic(err)
//	}
//	a := "hello?"
//	goodsClient := goods.NewGoodsServiceClient(client)
//	searchGoods, err := goodsClient.SearchGoods(context.Background(), &goods.SearchRequest{SearchContent: &a})
//	if err != nil {
//		return
//	}
//	spew.Dump(searchGoods)
//	select {}
//}
