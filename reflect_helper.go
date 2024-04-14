package dynmgrm

import (
	"github.com/iancoleman/strcase"
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

// getDBNameFromStructField returns the name of the field in the struct
func getDBNameFromStructField(tf reflect.StructField) string {
	gTag := newGormTag(tf.Tag)
	if gTag.Column != "" {
		return gTag.Column
	}
	return strcase.ToSnake(tf.Name)
}
