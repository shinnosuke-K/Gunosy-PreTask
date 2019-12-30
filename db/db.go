package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Open() (*gorm.DB, error) {
	return gorm.Open("postgres", os.Getenv("DATABASE_URL")+"sslmode=disable")
}
