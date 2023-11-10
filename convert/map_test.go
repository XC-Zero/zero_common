package convert

import (
	"log"
	"reflect"
	"testing"
)

func TestMapConvert(t *testing.T) {
	var a = make(map[string]any)
	var b = make(map[string]string)
	var c = make(map[string]struct{})

	var (
		ra = reflect.TypeOf(a)
		rb = reflect.TypeOf(b)
		rc = reflect.TypeOf(c)
	)

	log.Println(ra.ConvertibleTo(rb))
	log.Println(rb.ConvertibleTo(ra))
	log.Println(rc.ConvertibleTo(ra))
}
