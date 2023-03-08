package migrate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type tmplData struct {
	GormCliConfig
}

// CreateConnection - Create the migration file with template.
func CreateConnection(c GormCliConfig) (string, error) {
	// Create migration folder anyway.
	_ = os.Mkdir(MigrationTargetFolder, os.ModePerm)

	dialectString := ""
	dsn := ""

	data := tmplData{GormCliConfig: c}
	switch c.DB.Dialects {
	case "mysql":
		dialectString = `"gorm.io/driver/mysql"`
		dsn = `"{{.DB.User}}:{{.DB.Password}}@tcp({{.DB.Host}}:{{.DB.Port}})/{{.DB.Dbname}}?charset={{.DB.Charset}}&parseTime=True&loc=Local"`
	case "postgres":
		if data.DB.TimeZone == "" {
			data.DB.TimeZone = "UTC"
		}
		if data.DB.SSLMode == "" {
			data.DB.SSLMode = "disable"
		}
		dialectString = `"gorm.io/driver/postgres"`
		dsn = `"host={{.DB.Host}} user={{.DB.User}} password={{.DB.Password}} dbname={{.DB.Dbname}} port={{.DB.Port}} sslmode={{.DB.SSLMode}} TimeZone={{.DB.TimeZone}}"`
	}

	connectionTemplate := fmt.Sprintf(
		`package main

import (
	"gorm.io/gorm"
	%v
)

// NewDB - Get gorm DB instance.
func NewDB() (*gorm.DB, error) {
	dsn := %v
	db, err := gorm.Open({{.DB.Dialects}}.Open(dsn), &gorm.Config{})
	return db, err
}`, dialectString, dsn)

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
		return "", fmt.Errorf("%v (%v)", "Create connection failed.", err.Error())
	}

	if err := RemovePlugin(targetFileName); err != nil {
		return "", err
	}
	pluginFile, err := BuildPlugin(targetFileName)
	if err != nil {
		return "", err
	}

	if os.Getenv("DEBUG_CONNECTION") != "true" {
		if err := os.Remove(tmpFile); err != nil {
			return "", fmt.Errorf("%v (%v)", "Remove temp file failed", err.Error())
		}
	}

	return pluginFile, nil
}
