package convert

import (
	"log"
	"testing"
)

func TestAmount_Float64(t *testing.T) {
	var a = []string{"−502,65", "12,567,56", "-12 567.02", "78965.28"}
	for _, s := range a {
		log.Println(s, IsVailAmount(s))
		if IsVailAmount(s) {
			log.Println(NewAmount(s).SplitAll(".", ",", " ").Float64())
		}
	}

}

func TestToNormalAmount(t *testing.T) {
	log.Println(IsVailAmount("1 187,95"))
	//log.Println(ToNormalAmount("1.023,98"))
	log.Println(ToNormalAmount("1 187,95"))

}
