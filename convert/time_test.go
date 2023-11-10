package convert

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTryToLocalTIme(t *testing.T) {
	parse, err := time.Parse("2006/01/02 15:04:05 MST", "1998/12/02 12:00:55 DST")
	if err != nil {
		panic(err)
	}
	log.Println(parse, parse.IsDST())
	location, err := time.LoadLocation("CST")
	if err != nil {
		panic(err)
	}
	log.Println(location)
}

func TestMaybeTime(t *testing.T) {
	log.Println(MaybeTime("31.08.2023 23:55:48 UTC"))
}

func TestName(t *testing.T) {
	log.Println(MaybeTime(`9 apr. 2023 15:51:04 UTC`))
	transformTime, err := TransformTime(`2 Jan 2006 15:04:05 MST`, `9 apr. 2023 15:51:04 UTC`)
	if err != nil {
		panic(err)
	}
	location, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		panic(err)

	}
	log.Println(transformTime.In(location))

}

func TestTime(t *testing.T) {

	parse, err := time.Parse(time.RFC3339, "2023-06-01T18:00:00+00:00")
	if err != nil {
		panic(err)
	}
	location, err := time.LoadLocation("WET")
	if err != nil {
		panic(err)
	}
	log.Println(parse.In(location))
}

func TestAbandonTimeZone(t *testing.T) {
	parse, err := time.Parse(time.RFC3339, "2023-06-20T20:23:19-07:00")
	if err != nil {
		panic(err)
	}
	log.Println(AbandonTimeZone(parse))

	log.Println(fmt.Sprintf("%v", 1)[0] == '1')
}
