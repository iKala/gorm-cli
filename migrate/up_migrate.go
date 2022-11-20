package migrate

import (
	"fmt"
	"os"
	"strings"

	"errors"

	"gorm.io/gorm"
)

// UpMigrate - Run migration
func UpMigrate(db *gorm.DB, files []os.FileInfo) error {
	if len(files) == 0 {
		return ErrNoMigration
	}

	for i, file := range files {
		if file.Name() == "connection.so" || (!strings.HasSuffix(file.Name(), ".go") && !strings.HasSuffix(file.Name(), ".so")) {
			continue
		}
		var meta GormMeta
		err := db.Where(&GormMeta{Name: file.Name()}).First(&meta).Error

		// Execute the migration when the record not found.
		if err == nil {
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		pluginName, err := BuildPlugin(file.Name())
		if err != nil {
			return fmt.Errorf("%v (%v)", "Build plugin failed.", err.Error())
		}

		migration, err := getMigration(pluginName)
		if err != nil {
			return fmt.Errorf("%v (%v)", "Load migration plugin failed", err.Error())
		}

		if err := migration.Up(db); err != nil {
			return fmt.Errorf("%v (%v)", "Migrate failed."+file.Name(), err.Error())
		}

		fmt.Println("Migrated.", i, file.Name())
		if err := db.Create(&GormMeta{Name: file.Name()}).Error; err != nil {
			return err
		}
	}

	return nil
}
