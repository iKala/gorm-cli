package migrate

import (
	"gorm.io/gorm"
)

// GormCliConfig is the config for how gorm-cli connecting to DB
type GormCliConfig struct {
	DB struct {
		Dialects string
		Port     string
		Host     string
		User     string
		Password string
		Dbname   string
		Charset  string
	}
	Migration struct {
		Path string
	}
}

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
