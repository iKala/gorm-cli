package migrate

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// CreateMigration - Create the migration file with template.
func CreateMigration(purpose string) (string, error) {
	if purpose == "" {
		return "", ErrEmptyPurpose
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
		return "", ErrDuplicatedMigration
	}

	err := ioutil.WriteFile(
		targetFileName,
		[]byte(migrationTemplate),
		os.ModePerm,
	)

	if err != nil {
		return "", fmt.Errorf("%v (%v)", "Create migration failed.", err.Error())
	}

	return targetFileName, nil
}
