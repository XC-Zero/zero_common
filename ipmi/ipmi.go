package main

import (
	"github.com/bougou/go-ipmi"
	"github.com/pkg/errors"
	"log"
)

func Ipmi() {
	//client, err := ipmi.NewClient("192.168.12.60", 23, "tessan", "Tessan123456")
	client, err := ipmi.NewOpenClient()
	if err != nil {
		log.Panicf(`%+v`, errors.WithStack(err))
	}
	if err := client.Connect(); err != nil {
		log.Panicf(`%+v`, errors.WithStack(err))
	}
	sensors, err := client.GetSensors()
	if err != nil {
		log.Panicf(`%+v`, errors.WithStack(err))

	}

	for i := range sensors {
		switch sensors[i].SensorType {
		case ipmi.SensorTypeFan, ipmi.SensorTypeCurrent, ipmi.SensorTypeOtherFRU,
			ipmi.SensorTypeOtherUnitsbased, ipmi.SensorTypeTemperature, ipmi.SensorTypeVoltage:
			log.Printf(" %+v \n", sensors[i])
		}

	}

}

func main() {
	Ipmi()
}
