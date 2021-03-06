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
}

// CreateConnection - Create the migration file with template.
func CreateConnection(c GormCliConfig) (string, error) {
	// Create migration folder anyway.
	_ = os.Mkdir(MigrationTargetFolder, os.ModePerm)

	data := tmplData{GormCliConfig: c}
	if c.DB.Dialects == "mysql" {
		data.DialectString = `_ "github.com/jinzhu/gorm/dialects/mysql"`
	}

	connectionTemplate :=
		`package main

import (
	"github.com/jinzhu/gorm"
	{{.DialectString}}
)

// NewDB - Get gorm DB instance.
func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open("{{.DB.Dialects}}", "{{.DB.User}}:{{.DB.Password}}@tcp({{.DB.Host}})/{{.DB.Dbname}}?charset={{.DB.Charset}}&parseTime=True&loc=Local")
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
