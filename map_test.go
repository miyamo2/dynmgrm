package dynamgorm

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestMap_GormDataType(t *testing.T) {
	type testCase struct {
		sut  Map
		want string
	}
	tests := map[string]testCase{
		"happy-path": {
			sut:  Map{},
			want: "dgmap",
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

func TestMap_Scan(t *testing.T) {
	type testCase struct {
		sut           Map
		args          interface{}
		want          error
		expectedState Map
	}
	tests := map[string]testCase{
		"happy-path/empty-map": {
			sut:           Map{},
			args:          map[string]interface{}{},
			expectedState: Map{},
		},
		"happy-path/single-value": {
			sut:           Map{},
			args:          map[string]interface{}{"a": 1},
			expectedState: Map{"a": 1},
		},
		"happy-path/multiple-values": {
			sut:           Map{},
			args:          map[string]interface{}{"a": 1, "b": "2"},
			expectedState: Map{"a": 1, "b": "2"},
		},
		"unhappy-path/non-map-value": {
			sut:  Map{},
			args: "non-map",
			want: ErrFailedToCast,
		},
		"unhappy-path/sut-is-not-empty": {
			sut:           Map{"existing": 1},
			args:          map[string]interface{}{"a": 1},
			want:          ErrCollectionAlreadyContainsItem,
			expectedState: Map{"existing": 1},
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

func TestMap_ResolveNestedCollections(t *testing.T) {
	type testCase struct {
		sut           Map
		expectedState Map
		want          error
	}
	tests := map[string]testCase{
		"happy-path/empty-map": {
			sut:           Map{},
			expectedState: Map{},
		},
		"happy-path/single-nested-map": {
			sut:           Map{"a": map[string]interface{}{"b": 1}},
			expectedState: Map{"a": Map{"b": 1}},
		},
		"happy-path/multiple-nested-maps": {
			sut:           Map{"a": map[string]interface{}{"b": 1}, "c": map[string]interface{}{"d": 2}},
			expectedState: Map{"a": Map{"b": 1}, "c": Map{"d": 2}},
		},
		"happy-path/nested-list": {
			sut:           Map{"a": []interface{}{1, "b"}},
			expectedState: Map{"a": List{1, "b"}},
		},
		"happy-path/nested-int-sets": {
			sut:           Map{"a": []interface{}{int(1), int(2), int(3)}},
			expectedState: Map{"a": Sets[int]{1, 2, 3}},
		},
		"happy-path/nested-float-sets": {
			sut:           Map{"a": []interface{}{float64(1.1), float64(2.1), float64(3.1)}},
			expectedState: Map{"a": Sets[float64]{1.1, 2.1, 3.1}},
		},
		"happy-path/nested-string-sets": {
			sut:           Map{"a": []interface{}{string("1"), string("2"), string("3")}},
			expectedState: Map{"a": Sets[string]{"1", "2", "3"}},
		},
		"happy-path/nested-binary-sets": {
			sut:           Map{"a": []interface{}{[]byte("1"), []byte("2"), []byte("3")}},
			expectedState: Map{"a": Sets[[]byte]{[]byte("1"), []byte("2"), []byte("3")}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := resolveCollectionsNestedInMap(&tt.sut)
			if !errors.Is(err, tt.want) {
				t.Errorf("ResolveNestedDocument() error = %v, want %v", err, tt.want)
				return

			}
			if diff := cmp.Diff(tt.expectedState, tt.sut); diff != "" {
				t.Errorf("ResolveNestedDocument() mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
