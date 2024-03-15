package types

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSets_Scan(t *testing.T) {
	type testCase struct {
		sut           Sets[string]
		args          interface{}
		want          error
		expectedState Sets[string]
	}
	tests := map[string]testCase{
		"unhappy-path/already-contains-item": {
			sut:           Sets[string]{"value"},
			args:          []interface{}{"test"},
			expectedState: Sets[string]{"value"},
			want:          ErrCollectionAlreadyContainsItem,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestSets_Scan_String(t *testing.T) {
	type testCase struct {
		sut           Sets[string]
		args          interface{}
		want          error
		expectedState Sets[string]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           Sets[string]{},
			args:          []interface{}{"test"},
			expectedState: Sets[string]{"test"},
		},
		"happy-path/multiple-values": {
			sut:           Sets[string]{},
			args:          []interface{}{"test1", "test2"},
			expectedState: Sets[string]{"test1", "test2"},
		},
		"happy-path/numeric-value": {
			sut:           Sets[string]{},
			args:          []interface{}{1.1, 1.2},
			expectedState: Sets[string]{"1.1", "1.2"},
		},
		"happy-path/numeric-and-string": {
			sut:           Sets[string]{},
			args:          []interface{}{1.1, "1.2"},
			expectedState: Sets[string]{"1.1", "1.2"},
		},
		"happy-path/null": {
			sut:           Sets[string]{},
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           Sets[string]{},
			args:          "test",
			want:          ErrValueIsIncompatibleOfInterfaceSlice,
			expectedState: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestSets_Scan_Int(t *testing.T) {
	type testCase struct {
		sut           Sets[int]
		args          interface{}
		want          error
		expectedState Sets[int]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           Sets[int]{},
			args:          []interface{}{1},
			expectedState: Sets[int]{1},
		},
		"happy-path/multiple-values": {
			sut:           Sets[int]{},
			args:          []interface{}{1, 2},
			expectedState: Sets[int]{1, 2},
		},
		"happy-path/float-value": {
			sut:           Sets[int]{},
			args:          []interface{}{1.1, 1.2},
			expectedState: Sets[int]{1, 1},
		},
		"happy-path/int-and-float32": {
			sut:           Sets[int]{},
			args:          []interface{}{1, float32(1.2)},
			expectedState: Sets[int]{1, 1},
		},
		"happy-path/int-and-float64": {
			sut:           Sets[int]{},
			args:          []interface{}{1, float64(1.2)},
			expectedState: Sets[int]{1, 1},
		},
		"happy-path/float32-and-float64": {
			sut:           Sets[int]{},
			args:          []interface{}{float32(1.1), float64(1.2)},
			expectedState: Sets[int]{1, 1},
		},
		"happy-path/null": {
			sut:           Sets[int]{},
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           Sets[int]{},
			args:          1,
			want:          ErrValueIsIncompatibleOfInterfaceSlice,
			expectedState: nil,
		},
		"unhappy-path/not-compatible-slice": {
			sut:           Sets[int]{},
			args:          []interface{}{"A"},
			want:          ErrValueIsIncompatibleOfIntSlice,
			expectedState: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestSets_Scan_Float32(t *testing.T) {
	type testCase struct {
		sut           Sets[float32]
		args          interface{}
		want          error
		expectedState Sets[float32]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           Sets[float32]{},
			args:          []interface{}{float32(1.1)},
			expectedState: Sets[float32]{float32(1.1)},
		},
		"happy-path/multiple-values": {
			sut:           Sets[float32]{},
			args:          []interface{}{float32(1.1), float32(1.2)},
			expectedState: Sets[float32]{float32(1.1), float32(1.2)},
		},
		"happy-path/int-value": {
			sut:           Sets[float32]{},
			args:          []interface{}{1, 2},
			expectedState: Sets[float32]{1, 2},
		},
		"happy-path/int-and-float32": {
			sut:           Sets[float32]{},
			args:          []interface{}{1, float32(1.2)},
			expectedState: Sets[float32]{1, 1.2},
		},
		"happy-path/null": {
			sut:           Sets[float32]{},
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           Sets[float32]{},
			args:          1,
			want:          ErrValueIsIncompatibleOfInterfaceSlice,
			expectedState: nil,
		},
		"unhappy-path/not-compatible-slice": {
			sut:           Sets[float32]{},
			args:          []interface{}{"A"},
			want:          ErrValueIsIncompatibleOfFloat32Slice,
			expectedState: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestSets_Scan_Float64(t *testing.T) {
	type testCase struct {
		sut           Sets[float64]
		args          interface{}
		want          error
		expectedState Sets[float64]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           Sets[float64]{},
			args:          []interface{}{1.1},
			expectedState: Sets[float64]{1.1},
		},
		"happy-path/multiple-values": {
			sut:           Sets[float64]{},
			args:          []interface{}{1.1, 1.2},
			expectedState: Sets[float64]{1.1, 1.2},
		},
		"happy-path/int-value": {
			sut:           Sets[float64]{},
			args:          []interface{}{1, 2},
			expectedState: Sets[float64]{1, 2},
		},
		"happy-path/int-and-float64": {
			sut:           Sets[float64]{},
			args:          []interface{}{1, 1.2},
			expectedState: Sets[float64]{1, 1.2},
		},
		"happy-path/null": {
			sut:           Sets[float64]{},
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:  Sets[float64]{},
			args: 1,
			want: ErrValueIsIncompatibleOfInterfaceSlice,
		},
		"unhappy-path/not-compatible-slice": {
			sut:  Sets[float64]{},
			args: []interface{}{"A"},
			want: ErrValueIsIncompatibleOfFloat64Slice,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestSets_Scan_Binary(t *testing.T) {
	type testCase struct {
		sut           Sets[[]byte]
		args          interface{}
		want          error
		expectedState Sets[[]byte]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           Sets[[]byte]{},
			args:          []interface{}{[]byte{116, 101, 115, 116, 49}},
			expectedState: Sets[[]byte]{[]byte{116, 101, 115, 116, 49}},
		},
		"happy-path/multiple-value": {
			sut: Sets[[]byte]{},
			args: []interface{}{
				[]byte{116, 101, 115, 116, 49},
				[]byte{116, 101, 115, 116, 50},
			},
			expectedState: Sets[[]byte]{
				[]byte{116, 101, 115, 116, 49},
				[]byte{116, 101, 115, 116, 50},
			},
		},
		"happy-path/string-value": {
			sut:           Sets[[]byte]{},
			args:          []interface{}{"test"},
			expectedState: Sets[[]byte]{[]byte("test")},
		},
		"happy-path/null": {
			sut:           Sets[[]byte]{},
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:  Sets[[]byte]{},
			args: 1,
			want: ErrValueIsIncompatibleOfInterfaceSlice,
		},
		"unhappy-path/not-compatible-slice": {
			sut:  Sets[[]byte]{},
			args: []interface{}{1},
			want: ErrValueIsIncompatibleOfBinarySlice,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.sut.Scan(tt.args)
			if !errors.Is(err, tt.want) {
				t.Errorf("Scan() error = %v, want %v", err, tt.want)
				return
			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("Scan() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}

func TestSets_GormDataType(t *testing.T) {
	type testCase struct {
		sut  Sets[string]
		want string
	}
	tests := map[string]testCase{
		"happy-path": {
			sut:  Sets[string]{},
			want: "dgsets",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.sut.GormDataType(); got != tt.want {
				t.Errorf("GormDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}
