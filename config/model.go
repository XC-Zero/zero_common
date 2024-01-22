package config

//// ReadConfig 读取配置文件
//func ReadConfig(path string, config any) error {
//	file, err := os.ReadFile(path)
//	if err != nil {
//		return err
//	}
//	viper.SetConfigType("toml")
//	err = viper.ReadConfig(bytes.NewBuffer(file))
//	if err != nil {
//		return err
//	}
//	if reflect.TypeOf(config).Kind() != reflect.Ptr {
//		return errors.New("Config is not Ptr! ")
//	}
//	err = viper.Unmarshal(config, setTagName)
//	if err != nil {
//		return err
//	}
//	return nil
//}

type MongoDBConfig struct {
	URL    string `json:"url" yaml:"url"  toml:"url"`
	DBName string `json:"db_name" yaml:"db_name" toml:"db_name"`
}

type InfluxDBConfig struct {
	Host     string `json:"host" yaml:"host" toml:"host"`
	Token    string `json:"token" yaml:"token" toml:"token"`
	Org      string `json:"org" yaml:"org" toml:"org"`
	Database string `json:"database" yaml:"database" toml:"database"`
}

type MinioConfig struct {
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	EndPoint  string `json:"end_point" yaml:"end_point"`
	Bucket    string `json:"bucket" yaml:"bucket"`
}

type MysqlConfig struct {
	Host     string `json:"host" toml:"host"`
	Port     string `json:"port" toml:"port"`
	Username string `json:"username" toml:"username"`
	Pass     string `json:"pass" toml:"pass"`
	DBName   string `json:"db_name" toml:"db_name"`
	LogMode  int    `json:"log_mode" toml:"log_mode"`
	TZ       string `json:"timezone" toml:"timezone"`
}

type EmailConfig struct {
	EmailServerAddr string `json:"email_server_addr" yaml:"email_server_addr" toml:"email_server_addr"`
	SenderEmail     string `json:"sender_email" yaml:"sender_email" toml:"sender_email"`
	EmailSecret     string `json:"email_secret" yaml:"email_secret" toml:"email_secret"`
}

type MonitorConfig struct {
	CpuPercentThresholds  int      `json:"cpu_percent_thresholds"`
	MemPercentThresholds  int      `json:"mem_percent_thresholds"`
	DiskPercentThresholds int      `json:"disk_percent_thresholds"`
	CheckFrequency        string   `json:"check_frequency"`
	SendEmailToList       []string `json:"send_email_to_list"`
}

type KafkaConfig struct {
	Address []string `json:"address" toml:"address"`
	SSL     bool     `json:"ssl" toml:"ssl"`
}
type RedisConfig struct {
	AddrList []string `json:"addr_list" yaml:"addr_list"`
	Password string   `json:"password" yaml:"password"`
}

type ElasticSearchConfig struct {
	Host                         []string `json:"host" yaml:"host"`
	User                         string   `json:"user" yaml:"user"`
	Password                     string   `json:"password" yaml:"password"`
	ResponseHeaderTimeoutSeconds int      `json:"response_header_timeout_seconds" yaml:"response_header_timeout_seconds"`
}
