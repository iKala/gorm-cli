package migrate

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// UpMigrate - Run migration
func UpMigrate(db *gorm.DB, files []os.FileInfo) error {
	for i, file := range files {
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		var meta GormMeta
		err := db.Where(&GormMeta{Name: file.Name()}).First(&meta).Error

		// Execute the migration when the record not found.
		if err == nil {
			continue
		}

		pluginName, err := buildPlugin(file.Name())
		if err != nil {
			return errors.Wrap(err, "Build plugin failed.")
		}

		migration, err := getMigration(pluginName)
		if err != nil {
			return errors.Wrap(err, "Load migration plugin failed")
		}

		if err := migration.Up(db); err != nil {
			return errors.Wrap(err, "Migrate failed."+file.Name())
		}

		fmt.Println("Migrated.", i, file.Name())
		if err := db.Create(&GormMeta{Name: file.Name()}).Error; err != nil {
			return err
		}

		if err := removePlugin(file.Name()); err != nil {
			return errors.Wrap(err, "Remove plugin file failed.")
		}
	}

	return nil
}
