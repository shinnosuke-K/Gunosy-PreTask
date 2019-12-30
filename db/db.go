package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Open() (*gorm.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL") + "?sslmode=disable"
	return gorm.Open("postgres", dbUrl)
}
