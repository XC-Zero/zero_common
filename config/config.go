package config

import (
	"flag"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"sync"
)

var (
	config      Config
	once        sync.Once
	configName  = flag.String("c", "config", "读入config文件")
	configPath  = flag.String("p", "./configs", "config所在目录")
	globalViper = viper.New()
)

// GetConfig 单例获取配置
func GetConfig() Config {
	once.Do(func() {
		flag.Parse()
		globalViper.SetConfigName(*configName)
		globalViper.AutomaticEnv()
		globalViper.AddConfigPath(*configPath)
		globalViper.SetConfigType("toml")
		if err := globalViper.ReadInConfig(); err != nil {
			panic(err)
		}

		err := globalViper.Unmarshal(&config, setTagName)
		if err != nil {
			panic(err)
		}
	})
	return config
}

type Config struct {
	Database Database      `json:"database"`
	Service  ServiceConfig `json:"service"`
	Monitor  MonitorConfig `json:"monitor"`
}

type Database struct {
	MysqlConfig           MysqlConfig         `json:"mysql"`
	LocalMysqlWriteConfig MysqlConfig         `json:"local_mysql_write"`
	LocalMysqlReadConfig  MysqlConfig         `json:"local_mysql_read"`
	MinioConfig           MinioConfig         `json:"minio"`
	KafkaConfig           KafkaConfig         `json:"kafka"`
	RedisConfig           RedisConfig         `json:"redis"`
	MongoConfig           MongoDBConfig       `json:"mongo"`
	ElasticSearchConfig   ElasticSearchConfig `json:"elasticsearch"`
	InfluxDBConfig        InfluxDBConfig      `json:"influxdb"`
}

type ServiceConfig struct {
	IsPro         bool          `json:"is_pro"`
	Port          string        `json:"port"`
	Debug         bool          `json:"debug"`
	AutoMigrate   bool          `json:"auto_migrate"`
	EmailConfig   EmailConfig   `json:"email"`
	PlumberConfig PlumberConfig `json:"plumber"`
}

type SyncMode int

const (
	SINGLE = iota
	BATCH
	INIT = 9
)

type PlumberConfig struct {
	ConsumerSuffix string   `json:"consumer_suffix"`
	CheckConnector bool     `json:"check_connector"`
	SyncMode       SyncMode `json:"sync_mode"`
}

// 设置 config 对应的结构体的 tag
func setTagName(d *mapstructure.DecoderConfig) {
	d.TagName = "json"
}
