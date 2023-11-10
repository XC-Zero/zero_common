package sqlite3

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
)

func InitClient(path string) (*gorm.DB, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		create, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		create.Close()

	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}
