package migrate

import (
	"errors"
)

var (
	ErrMigrationCanceled   = errors.New("rollback migration canceled")
	ErrEmptyPurpose        = errors.New("missing purpose when creating migration")
	ErrDuplicatedMigration = errors.New("migration exists")
	ErrNoMigration         = errors.New("running with no migration")
)
