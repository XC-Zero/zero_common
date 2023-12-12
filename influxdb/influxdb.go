package influxdb

import (
	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/XC-Zero/zero_common/config"
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
