package migrate

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"gorm.io/gorm"
)

// DownMigration - Rollback migration
func DownMigration(db *gorm.DB, files []os.FileInfo, step int64) error {
	if len(files) == 0 {
		return ErrNoMigration
	}

	stepMessage := strconv.Itoa(int(step))
	if step == -1 {
		stepMessage = "all"
	}
	prompt := promptui.Prompt{
		Label: "Are you sure to rollback the migration with " + stepMessage + " steps? (Yes/No)",
	}

	result, err := prompt.Run()

	if err != nil || result != "Yes" {
		return ErrMigrationCanceled
	}

	var metas []GormMeta
	if err := db.Order("ID desc").Limit(int(step)).Find(&metas).Error; err != nil {
		return err
	}

	for i, meta := range metas {
		pluginName, err := BuildPlugin(meta.Name)
		if err != nil {
			fmt.Println("Build plugin failed.", err)
			return err
		}

		migration, err := getMigration(pluginName)
		if err != nil {
			return err
		}

		if err := migration.Down(db); err != nil {
			return err
		}

		fmt.Println("Rollbacked.", i, meta.Name)
		if err := db.Delete(&meta).Error; err != nil {
			return err
		}
	}

	return nil
}
