package dynmgrm

import (
	"errors"

	"github.com/btnguyen2k/godynamo"

	"gorm.io/gorm"
)

// Translate it will translate the error to native gorm errors.
func (dialector Dialector) Translate(err error) error {
	switch {
	case errors.Is(err, godynamo.ErrTxCommitting),
		errors.Is(err, godynamo.ErrTxRollingBack),
		errors.Is(err, godynamo.ErrInTx),
		errors.Is(err, godynamo.ErrInvalidTxStage),
		errors.Is(err, godynamo.ErrNoTx):
		return gorm.ErrInvalidTransaction
	}
	return err
}
