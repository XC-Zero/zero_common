package main

import (
	"context"
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/XC-Zero/zero_common/config"
	"github.com/XC-Zero/zero_common/influxdb"
	"github.com/bougou/go-ipmi"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"
)

var once sync.Once
var ipmiClient *ipmi.Client

func main() {
	saveData()

	ticker := time.NewTicker(time.Second * 5)
	for range ticker.C {
		saveData()
	}
	select {}
}

func getData() []*influxdb3.Point {
	once.Do(func() {
		client, err := ipmi.NewOpenClient()
		if err != nil {
			log.Panicf(`%+v`, errors.WithStack(err))
		}
		if err := client.Connect(); err != nil {
			log.Panicf(`%+v`, errors.WithStack(err))
		}
		ipmiClient = client
	})
	sensors, err := ipmiClient.GetSensors()
	if err != nil {
		log.Panicf(`%+v`, errors.WithStack(err))

	}
	var points []*influxdb3.Point
	for i := range sensors {
		switch sensors[i].SensorType {
		case ipmi.SensorTypeFan, ipmi.SensorTypeCurrent, ipmi.SensorTypeOtherFRU,
			ipmi.SensorTypeOtherUnitsbased, ipmi.SensorTypeTemperature, ipmi.SensorTypeVoltage:
			now := time.Now()
			points = append(points, influxdb3.NewPoint(sensors[i].Name, map[string]string{
				"unit": sensors[i].SensorUnit.String(),
			}, map[string]any{
				"value": sensors[i].Value,
			}, now))

		}

	}
	return points
}

func saveData() {
	ic := GetInstance()
	err := ic.WritePoints(context.Background(), getData()...)
	if err != nil {
		panic(err)
	}
}

var influxOnce sync.Once
var client *influxdb3.Client

func GetInstance() *influxdb3.Client {
	influxOnce.Do(func() {
		influxClient, err := influxdb.InitInfluxClient(config.GetConfig().Database.InfluxDBConfig)
		if err != nil {
			panic(err)
		}
		client = influxClient
	})
	return client
}
