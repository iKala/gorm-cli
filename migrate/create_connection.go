package migrate

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

type tmplData struct {
	GormCliConfig
	DialectString string
	Port          string
}

// CreateConnection - Create the migration file with template.
func CreateConnection(c GormCliConfig) (string, error) {
	// Create migration folder anyway.
	_ = os.Mkdir(MigrationTargetFolder, os.ModePerm)

	data := tmplData{GormCliConfig: c}
	if c.DB.Dialects == "mysql" {
		data.DialectString = `"gorm.io/driver/mysql"`
	}

	if c.DB.Port == "" {
		data.Port = "3306"
	}

	connectionTemplate :=
		`package main

import (
	"gorm.io/gorm"
	{{.DialectString}}
)

// NewDB - Get gorm DB instance.
func NewDB() (*gorm.DB, error) {
	dsn := "{{.DB.User}}:{{.DB.Password}}@tcp({{.DB.Host}}:{{.Port}})/{{.DB.Dbname}}?charset={{.DB.Charset}}&parseTime=True&loc=Local"
	db, err := gorm.Open({{.DB.Dialects}}.Open(dsn), &gorm.Config{})
	return db, err
}`

	tmpl, err := template.New("connection").Parse(connectionTemplate)
	if err != nil {
		return "", err
	}
	var connectionFileStringBuffer bytes.Buffer
	if err := tmpl.Execute(&connectionFileStringBuffer, &data); err != nil {
		return "", err
	}

	targetFileName := "connection.go"
	tmpFile := MigrationTargetFolder + "/" + targetFileName

	if err := ioutil.WriteFile(
		tmpFile,
		connectionFileStringBuffer.Bytes(),
		os.ModePerm,
	); err != nil {
		return "", errors.Wrap(err, "Create connection failed.")
	}

	if err := removePlugin(targetFileName); err != nil {
		return "", err
	}
	pluginFile, err := buildPlugin(targetFileName)
	if err != nil {
		return "", err
	}

	if err := os.Remove(tmpFile); err != nil {
		return "", errors.Wrap(err, "Remove temp file failed")
	}

	return pluginFile, nil
}
