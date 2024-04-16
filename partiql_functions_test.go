package dynmgrm

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"testing"
)

func Test_listAppend_expression(t *testing.T) {
	type args struct {
		db     *gorm.DB
		column string
	}
	type want struct {
		xp  string
		err error
	}
	type test struct {
		args args
		want want
	}
	tests := map[string]test{
		"happy_path": {
			args: args{
				db: &gorm.DB{
					Config: &gorm.Config{},
				},
				column: "A",
			},
			want: want{
				xp: "list_append(A, ",
			},
		},
		"unhappy_path": {
			args: args{
				db: &gorm.DB{
					Config: &gorm.Config{},
				},
				column: "A==true",
			},
			want: want{
				err: ErrInvalidColumnName,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			la := ListAppend()
			got := la.expression(tt.args.db, tt.args.column)
			err := tt.args.db.Error
			if !errors.Is(err, tt.want.err) {
				t.Errorf("expression() error = %v, want nil", err)
			}
			if diff := cmp.Diff(tt.want.xp, got); diff != "" {
				t.Errorf("expression() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
