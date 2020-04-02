package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/iKala/gorm-cli/migrate"
	"gopkg.in/yaml.v2"
)

func main() {
	migrate.MigrationTargetFolder = "./migration"

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Something wrong when getting current path.")
		return
	}

	if len(os.Args) == 1 {
		fmt.Println("Option of migration is needed - (db:init db:migrate / db:rollback / db:create_migration)")
		return
	}

	var migrateAction string
	if migrateAction = os.Args[1]; migrateAction == "" {
		fmt.Println("Empty option is not allowed - (db:init db:migrate / db:rollback / db:create_migration)")
		return
	}
	if migrateAction == "db:init" {
		bytes, err := ioutil.ReadFile(path + "/.gorm-cli.yaml")
		if err != nil {
			fmt.Println(".gorm-cli.yaml not exists, you must create one for gorm-cli connecting to DB. https://github.com/iKala/gorm-cli/blob/master/README.md")
			return
		}

		connection := migrate.Connection{}
		if err := yaml.Unmarshal(bytes, &connection); err != nil {
			fmt.Println("Failed to parse .gorm-cli.yaml, might be syntax error. https://github.com/iKala/gorm-cli/blob/master/.gorm-cli.yaml")
			return
		}

		fileName, err := migrate.CreateConnection(connection)
		if err != nil {
			fmt.Println("Initiail connection file failed.")
			return
		}

		fmt.Println("Connection file created", fileName)
		return
	}

	rollbackStep := int64(-1)
	if migrateAction == "db:rollback" && len(os.Args) == 3 {
		argStep, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			fmt.Println("Wrong option format at args[2] - Rollback step needs to be a integer number.", err)
			return
		}

		rollbackStep = argStep
	}
	var purpose string
	if migrateAction == "db:create_migration" {
		if len(os.Args) == 3 {
			purpose = os.Args[2]
		} else {
			fmt.Println("Missing option at args[2] - The purpose must not empty with db:create_migration.")
			return
		}
	}

	files, err := ioutil.ReadDir(migrate.MigrationTargetFolder)
	if err != nil {
		fmt.Println("Migration folder not exists.")
		return
	}

	db := migrate.NewDB()
	defer db.Close()

	db.AutoMigrate(&migrate.GormMeta{})

	switch migrateAction {
	case "db:migrate":
		if err := migrate.UpMigrate(db, files); err != nil {
			fmt.Println(err)
		}
	case "db:rollback":
		if err := migrate.DownMigration(db, files, rollbackStep); err != nil {
			fmt.Println(err)
		}
	case "db:create_migration":
		fileName, err := migrate.CreateMigration(purpose)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Migration created.", fileName)
	}
}
