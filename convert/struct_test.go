package convert

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"testing"
)

type A struct {
	A string
}
type B struct {
	A string
}

func TestMapperStructByName(t *testing.T) {
	var a = A{
		A: "hello?",
	}
	var b B
	MapperStructByName(&b, &a)
	spew.Dump(b)
}

type ADC struct {
	Distance int
	Name     string
	cd       int
}

func TestCompareDifference(t *testing.T) {
	var jin = ADC{
		Distance: 525,
		Name:     "jin",
		cd:       4,
	}
	var nvqiang = ADC{
		Distance: 525,
		Name:     "nvqiang",
		cd:       5,
	}

	difference, err := CompareDifference(jin, nvqiang, "a", "b", "json", "Distance")
	if err != nil {
		panic(err)
	}
	spew.Dump(difference)

}

type WTF struct {
	A int `json:"a"`
}

func TestInvisibleDataCol(t *testing.T) {
	var wtf = WTF{
		A: 15,
	}
	aa := InvisibleDataCol([]WTF{wtf}, "a")
	log.Println(aa)
}
