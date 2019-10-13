package migrate

import (
	"github.com/jinzhu/gorm"
)

// Migration - The main interface for migration files
type Migration interface {
	Up(*gorm.DB) error
	Down(*gorm.DB) error
}

// GormMeta - The meta for storing which migration has been executed.
type GormMeta struct {
	gorm.Model

	Name string `gorm:"type:varchar(255);primary_key"`
}
