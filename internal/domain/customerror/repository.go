package customerror

import (
	"errors"
)

var (
	ErrEntityIsNull   = errors.New("entity is nil")
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrGormScan       = errors.New("gorm scanning customerror")
)
