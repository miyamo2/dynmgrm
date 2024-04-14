package dynmgrm

import (
	"reflect"
	"strings"
)

type gormTag struct {
	Column string
}

// gormTagWithStructTag returns gormTag representation of reflect.StructTag
func newGormTag(tag reflect.StructTag) gormTag {
	gTag := gormTag{}
	for _, value := range strings.Split(tag.Get("gorm"), ";") {
		if value == "" {
			continue
		}
		kv := strings.Split(value, ":")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "column":
			gTag.Column = kv[1]
		}
	}
	return gTag
}
