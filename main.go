package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/iKala/gorm-cli/migrate"
	"gopkg.in/yaml.v2"
)

func main() {
	c := migrate.GormCliConfig{}
	c.DB.Dialects = os.Getenv("DB_DIALECTS")
	c.DB.Dbname = os.Getenv("DB_DBNAME")
	c.DB.Host = os.Getenv("DB_HOST")
	c.DB.User = os.Getenv("DB_USER")
	c.DB.Password = os.Getenv("DB_PASSWORD")
	c.DB.Port = os.Getenv("DB_PORT")
	c.DB.Charset = os.Getenv("DB_CHARTSET")
	c.Migration.Path = os.Getenv("MIGRATION_PATH")

	migrate.MigrationTargetFolder = c.Migration.Path

	if len(os.Args) == 1 {
		fmt.Println("Option of migration is needed - (db:prebuild / db:init / db:migrate / db:rollback / db:create_migration)")
		return
	}

	var migrateAction string
	if migrateAction = os.Args[1]; migrateAction == "" {
		fmt.Println("Empty option is not allowed - (db:prebuild / db:init / db:migrate / db:rollback / db:create_migration)")
		return
	}
	if migrateAction == "db:prebuild" {
		fileInfo, err := ioutil.ReadDir(migrate.MigrationTargetFolder)
		if err != nil {
			fmt.Println("Failed to read migrations", err)
			return
		}
		for _, f := range fileInfo {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".go") {
				if err := migrate.RemovePlugin(f.Name()); err != nil {
					fmt.Println("Failed to reset existed plugin", f.Name(), err)
					return
				}

				createdFile, err := migrate.BuildPlugin(f.Name())
				if err != nil {
					fmt.Println("Failed to build plugin", f.Name(), err)
					return
				}

				fmt.Println(createdFile, "created")
			}
		}
	}

	if migrateAction == "db:init" {
		// Parse yaml setting when got non of env
		if reflect.DeepEqual(c, migrate.GormCliConfig{}) {
			path, err := os.Getwd()
			if err != nil {
				fmt.Println("Something wrong when getting current path.")
				return
			}

			bytes, err := ioutil.ReadFile(path + "/.gorm-cli.yaml")
			if err != nil {
				fmt.Println(".gorm-cli.yaml not exists, you must create one for gorm-cli connecting to DB. https://github.com/iKala/gorm-cli/blob/master/README.md")
				return
			}
			if err := yaml.Unmarshal(bytes, &c); err != nil {
				fmt.Println("Failed to parse .gorm-cli.yaml, might be syntax error. https://github.com/iKala/gorm-cli/blob/master/.gorm-cli.yaml")
				return
			}
		}
		fileName, err := migrate.CreateConnection(c)
		if err != nil {
			fmt.Println("Initiail connection file failed.", err)
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

	// Read prebuilt plugins first.
	files, _ := ioutil.ReadDir(migrate.MigrationTargetFolder + "/.plugins")
	if len(files) == 0 {
		// Load go file when prebuilt plugins not exists.
		files, err = ioutil.ReadDir(migrate.MigrationTargetFolder)
		if err != nil {
			fmt.Println("Migration folder not exists.")
			return
		}
	}

	db := migrate.NewDB()

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
	default:
		fmt.Println("No matched action. (db:prebuild / db:init / db:migrate / db:rollback / db:create_migration)")
	}
}
