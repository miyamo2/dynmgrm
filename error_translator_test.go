package dynmgrm

import (
	"errors"
	"testing"

	"github.com/btnguyen2k/godynamo"
	"gorm.io/gorm"
)

func TestDialector_Translate(t *testing.T) {
	type test struct {
		args error
		want error
	}
	errOther := errors.New("other")
	tests := map[string]test{
		"happy_path/ErrTxCommitting": {
			args: godynamo.ErrInTx,
			want: gorm.ErrInvalidTransaction,
		},
		"happy_path/ErrTxRollingBack": {
			args: godynamo.ErrTxRollingBack,
			want: gorm.ErrInvalidTransaction,
		},
		"happy_path/ErrInTx": {
			args: godynamo.ErrInTx,
			want: gorm.ErrInvalidTransaction,
		},
		"happy_path/ErrInvalidTxStage": {
			args: godynamo.ErrInvalidTxStage,
			want: gorm.ErrInvalidTransaction,
		},
		"happy_path/ErrNoTx": {
			args: godynamo.ErrNoTx,
			want: gorm.ErrInvalidTransaction,
		},
		"happy_path/other": {
			args: errOther,
			want: errOther,
		},
		"happy_path/nil": {
			args: nil,
			want: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dialector := Dialector{}
			err := dialector.Translate(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Translate() error = %v, want %v", err, tt.want)
			}
		})
	}
}
