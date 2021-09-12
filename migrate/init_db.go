package migrate

import (
	"gorm.io/gorm"
)

// NewDB get gorm DB instance.
func NewDB() *gorm.DB {
	s, err := getPlugin("connection.so", "NewDB")
	if err != nil {
		panic(err)
	}

	db, err := s.(func() (*gorm.DB, error))()
	if err != nil {
		panic(err)
	}

	return db
}
