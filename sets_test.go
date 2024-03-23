package dynmgrm

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
			sut:           newSets[string](),
			args:          []string{"test"},
			expectedState: Sets[string]{"test"},
		},
		"happy-path/multiple-values": {
			sut:           newSets[string](),
			args:          []string{"test1", "test2"},
			expectedState: Sets[string]{"test1", "test2"},
		},
		"happy-path/numeric-value": {
			sut:           newSets[string](),
			args:          []string{"1.1", "1.2"},
			expectedState: Sets[string]{"1.1", "1.2"},
		},
		"happy-path/null": {
			sut:           newSets[string](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           newSets[string](),
			args:          "test",
			want:          ErrValueIsIncompatibleOfStringSlice,
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
			sut:           newSets[int](),
			args:          []float64{1},
			expectedState: Sets[int]{1},
		},
		"happy-path/multiple-values": {
			sut:           newSets[int](),
			args:          []float64{1, 2},
			expectedState: Sets[int]{1, 2},
		},
		"happy-path/float-value": {
			sut:           newSets[int](),
			args:          []float64{1.0, 2.0},
			expectedState: Sets[int]{1, 2},
		},
		"happy-path/null": {
			sut:           newSets[int](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:           newSets[int](),
			args:          1,
			want:          ErrValueIsIncompatibleOfIntSlice,
			expectedState: nil,
		},
		"unhappy-path/string-slice": {
			sut:           newSets[int](),
			args:          []string{"A"},
			want:          ErrValueIsIncompatibleOfIntSlice,
			expectedState: nil,
		},
		"unhappy-path/float-value": {
			sut:           newSets[int](),
			args:          []float64{1.1, 2.1},
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

func TestSets_Scan_Float64(t *testing.T) {
	type testCase struct {
		sut           Sets[float64]
		args          interface{}
		want          error
		expectedState Sets[float64]
	}
	tests := map[string]testCase{
		"happy-path/single-value": {
			sut:           newSets[float64](),
			args:          []float64{1.1},
			expectedState: Sets[float64]{1.1},
		},
		"happy-path/multiple-values": {
			sut:           newSets[float64](),
			args:          []float64{1.1, 1.2},
			expectedState: Sets[float64]{1.1, 1.2},
		},
		"happy-path/int-value": {
			sut:           newSets[float64](),
			args:          []float64{1, 2},
			expectedState: Sets[float64]{1, 2},
		},
		"happy-path/int-and-float64": {
			sut:           newSets[float64](),
			args:          []float64{1, 1.2},
			expectedState: Sets[float64]{1, 1.2},
		},
		"happy-path/null": {
			sut:           newSets[float64](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:  newSets[float64](),
			args: 1,
			want: ErrValueIsIncompatibleOfFloat64Slice,
		},
		"unhappy-path/not-compatible-slice": {
			sut:  newSets[float64](),
			args: []string{"A"},
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
			sut:           newSets[[]byte](),
			args:          [][]byte{{116, 101, 115, 116, 49}},
			expectedState: Sets[[]byte]{[]byte{116, 101, 115, 116, 49}},
		},
		"happy-path/multiple-value": {
			sut: newSets[[]byte](),
			args: [][]byte{
				{116, 101, 115, 116, 49},
				{116, 101, 115, 116, 50},
			},
			expectedState: Sets[[]byte]{
				[]byte{116, 101, 115, 116, 49},
				[]byte{116, 101, 115, 116, 50},
			},
		},
		"happy-path/null": {
			sut:           newSets[[]byte](),
			args:          nil,
			expectedState: nil,
		},
		"unhappy-path/not-slice": {
			sut:  newSets[[]byte](),
			args: 1,
			want: ErrValueIsIncompatibleOfBinarySlice,
		},
		"unhappy-path/not-compatible-slice": {
			sut:  newSets[[]byte](),
			args: []interface{}{1},
			want: ErrValueIsIncompatibleOfBinarySlice,
		},
		"unhappy-path/string-value": {
			sut:           newSets[[]byte](),
			args:          []interface{}{"test"},
			expectedState: nil,
			want:          ErrValueIsIncompatibleOfBinarySlice,
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
			sut:  newSets[string](),
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

func TestSets_IsCompatibleWithSets_String(t *testing.T) {
	type testCase struct {
		sut  Sets[string]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: []string{"test"},
			want: true,
		},
		"happy-path/int-value": {
			args: []int{1},
			want: false,
		},
		"happy-path/float64-value": {
			args: []float64{1.1},
			want: false,
		},
		"happy-path/float32-value": {
			args: []float32{1.1},
			want: false,
		},
		"happy-path/binary-value": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSets[string](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSets_IsCompatibleWithSets_Int(t *testing.T) {
	type testCase struct {
		sut  Sets[int]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: []int{1},
			want: true,
		},
		"happy-path/string-value": {
			args: []string{"test"},
			want: false,
		},
		"happy-path/float64-value": {
			args: []float64{1.1},
			want: false,
		},
		"happy-path/float64-typed-int-value": {
			args: []float64{1},
			want: true,
		},
		"happy-path/binary-value": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSets[int](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSets_IsCompatibleWithSets_Float64(t *testing.T) {
	type testCase struct {
		sut  Sets[float64]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: []float64{1.1},
			want: true,
		},
		"unhappy-path/string-value": {
			args: []string{"test"},
			want: false,
		},
		"unhappy-path/int-value": {
			args: []int{1},
			want: false,
		},
		"unhappy-path/float32-value": {
			args: []float32{1.1},
			want: false,
		},
		"unhappy-path/binary-value": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSets[float64](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSets_IsCompatibleWithSets_Binary(t *testing.T) {
	type testCase struct {
		sut  Sets[[]byte]
		args interface{}
		want bool
	}
	tests := map[string]testCase{
		"happy-path": {
			args: [][]byte{{116, 101, 115, 116, 49}},
			want: true,
		},
		"unhappy-path/string-value": {
			args: []string{"test"},
			want: false,
		},
		"happy-path/int-value": {
			args: []int{1},
			want: false,
		},
		"happy-path/float64-value": {
			args: []float64{1.1},
			want: false,
		},
		"happy-path/float32-value": {
			args: []float32{1.1},
			want: false,
		},
		"unhappy-path/not-slice": {
			args: 1,
			want: false,
		},
		"unhappy-path/nil": {
			args: nil,
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := isCompatibleWithSets[[]byte](tt.args)
			if got != tt.want {
				t.Errorf("isCompatibleWithSets() = %v, want %v", got, tt.want)
			}
		})
	}
}
