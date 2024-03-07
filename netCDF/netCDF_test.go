package common

import (
	"github.com/XC-Zero/zero_common/config"
	"github.com/XC-Zero/zero_common/mongo"
	"github.com/pkg/errors"
	"log"
	"testing"
)

func TestNetCDF_SaveMongo(t *testing.T) {
	client, err := mongo.InitMongoClient(config.MongoDBConfig{
		URL:    "mongodb://root:root123456@192.168.185.97:27017",
		DBName: "know_weather",
	})
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}

	c, err := NewCDF("C:\\Users\\XC\\workspace\\temp\\WCS001005_2023-12-21_00-00-13_dbs_682_50mTP.nc")
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	err = c.SaveMongo(client.Database("know_weather"), "test01x")
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
}
