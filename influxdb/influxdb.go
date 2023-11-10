package influxdb

import (
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"gitlab.tessan.com/data-center/tessan-erp-common/config"
)

func InitInfluxClient(config config.InfluxDBConfig) (*influxdb3.Client, error) {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         config.Host,
		Token:        config.Token,
		Database:     config.Database,
		Organization: config.Org,
	})
	return client, err

}
