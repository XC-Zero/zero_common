package mysql

import (
	"fmt"
	"github.com/XC-Zero/zero_common/config"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

func InitClient(config config.MysqlConfig) (*gorm.DB, error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	if strings.TrimSpace(config.TZ) == "" {
		config.TZ = "UTC"
	}
	config.TZ = strings.ReplaceAll(config.TZ, "/", "%2F")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		config.Username, config.Pass, config.Host, config.Port, config.DBName, config.TZ,
	)
	open, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(config.LogMode)),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	db, err := open.DB()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Second * 60)
	db.SetConnMaxLifetime(time.Minute * 15)

	return open, nil
}
