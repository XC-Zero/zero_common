package main

import (
	"log"
	"reflect"
)

type MyStruct struct {
	MyInt int
}

func main() {

	var c = []MyStruct{{MyInt: 5}}
	reflect.ValueOf(c).Index(0).FieldByName("MyInt").SetZero()
	log.Println(c)

}
