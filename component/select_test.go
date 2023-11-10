package component

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"reflect"
	"testing"
)

type a struct {
	A string `json:"a" dc:""`
	B int    `json:"b"`
}

func TestToDefaultCol(t *testing.T) {
	var av = []any{a{}}
	var model any
	ll := reflect.TypeOf(av).Kind()
	log.Println(ll)
	switch ll {
	case reflect.Struct:
		model = av
	case reflect.Array, reflect.Slice:
		bb := reflect.ValueOf(av).Index(0).Interface()
		log.Println(reflect.TypeOf(bb))
		model = bb
	}
	spew.Dump(ToDefaultCol(model))
}
