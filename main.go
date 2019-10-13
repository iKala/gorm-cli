package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/iKala/gorm-cli/migrate"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Option of migration is needed - (db:migrate / db:rollback / db:create_migration)")
		return
	}

	var migrateAction string
	if migrateAction = os.Args[1]; migrateAction == "" {
		fmt.Println("Empty option is not allowed - (db:migrate / db:rollback / db:create_migration)")
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

	migrate.MigrationTargetFolder = "./migration"

	files, err := ioutil.ReadDir(migrate.MigrationTargetFolder)
	if err != nil {
		fmt.Println("Migration folder not exists.")
		return
	}

	db := model.InitDB()
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
