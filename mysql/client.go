package mysql

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.tessan.com/data-center/tessan-erp-common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitClient(config config.MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Pass, config.Host, config.Port, config.DBName,
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
