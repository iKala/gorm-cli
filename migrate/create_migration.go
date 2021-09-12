package migrate

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
)

// CreateMigration - Create the migration file with template.
func CreateMigration(purpose string) (string, error) {
	if purpose == "" {
		return "", errors.New("Missing purpose when creating migration")
	}

	// Create migration folder anyway.
	_ = os.Mkdir("./migration", os.ModePerm)

	migrationTemplate :=
		`package main

import (
	"gorm.io/gorm"
)

type migration string

// Up - Changes for the migration.
func (m migration) Up(db *gorm.DB) error {
	return nil
}

// Down - Rollback changes for the migration.
func (m migration) Down(db *gorm.DB) error {
	return nil
}

var Migration migration`

	targetFileName := MigrationTargetFolder + "/" + time.Now().Format("20060102150405") + "_" + purpose + ".go"

	if fileExists(targetFileName) {
		return "", errors.New("Migration exists")
	}

	err := ioutil.WriteFile(
		targetFileName,
		[]byte(migrationTemplate),
		os.ModePerm,
	)

	if err != nil {
		return "", errors.Wrap(err, "Create migration failed.")
	}

	return targetFileName, nil
}
