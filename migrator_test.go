package dynmgrm

import "testing"

func TestTableClass_String(t *testing.T) {
	type test struct {
		sut  TableClass
		want string
	}
	tests := map[string]test{
		"happy_path/standard": {
			sut:  TableClassStandard,
			want: "STANDARD",
		},
		"happy_path/standard_ia": {
			sut:  TableClassStandardIA,
			want: "STANDARD_IA",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.sut.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
