package dynmgrm

import (
	"github.com/iancoleman/strcase"
	"reflect"
	"strings"
)

type gormTag struct {
	Column string
	DBType string
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
		case "type":
			gTag.DBType = kv[1]
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

// getDBTypeFromStructField returns the type of the field in the struct
func getDBTypeFromStructField(tf reflect.StructField) string {
	gTag := newGormTag(tf.Tag)
	return gTag.DBType
}

type secondaryIndexKind int

const (
	secondaryIndexKindLSI secondaryIndexKind = iota + 1
	secondaryIndexKindGSI
)

func (sik secondaryIndexKind) String() string {
	switch sik {
	case secondaryIndexKindLSI:
		return "LSI"
	case secondaryIndexKindGSI:
		return "GSI"
	}
	return ""
}

type secondaryIndexProperty struct {
	pk   bool
	sk   bool
	name string
	kind secondaryIndexKind
}

type dynmgrmTag struct {
	pk            bool
	sk            bool
	indexProperty []secondaryIndexProperty
	nonProjective []string
}

func newDynmgrmTag(tag reflect.StructTag) dynmgrmTag {
	res := dynmgrmTag{}
	for _, value := range strings.Split(tag.Get("dynmgrm"), ";") {
		if value == "" {
			continue
		}
		kv := strings.Split(value, ":")
		switch tn := kv[0]; tn {
		case "pk":
			if res.sk {
				continue
			}
			res.pk = true
		case "sk":
			if res.pk {
				continue
			}
			res.sk = true
		case "gsi-pk":
			iprp := secondaryIndexProperty{
				pk:   true,
				name: kv[1],
				kind: secondaryIndexKindGSI,
			}
			res.indexProperty = append(res.indexProperty, iprp)
		case "gsi-sk":
			iprp := secondaryIndexProperty{
				sk:   true,
				name: kv[1],
				kind: secondaryIndexKindGSI,
			}
			res.indexProperty = append(res.indexProperty, iprp)
		case "lsi":
			iprp := secondaryIndexProperty{
				name: kv[1],
				sk:   true,
				kind: secondaryIndexKindLSI,
			}
			res.indexProperty = append(res.indexProperty, iprp)
		case "non-projective":
			npl := strings.ReplaceAll(strings.ReplaceAll(kv[1], "[", ""), "]", "")
			for _, np := range strings.Split(npl, ",") {
				res.nonProjective = append(res.nonProjective, np)
			}
		}
	}
	return res
}

type dynmgrmKeyDefine struct {
	name     string
	dataType string
}

type dynmgrmSecondaryIndexDefine struct {
	pk                 dynmgrmKeyDefine
	sk                 dynmgrmKeyDefine
	nonProjectiveAttrs []string
}

type dynmgrmTableDefine struct {
	pk     dynmgrmKeyDefine
	sk     dynmgrmKeyDefine
	nonKey []string
	gsi    map[string]*dynmgrmSecondaryIndexDefine
	lsi    map[string]*dynmgrmSecondaryIndexDefine
}

func newDynmgrmTableDefine(modelMeta reflect.Type) dynmgrmTableDefine {
	res := dynmgrmTableDefine{
		nonKey: make([]string, 0, modelMeta.NumField()),
		gsi:    make(map[string]*dynmgrmSecondaryIndexDefine),
		lsi:    make(map[string]*dynmgrmSecondaryIndexDefine),
	}
	for i := 0; i < modelMeta.NumField(); i++ {
		tf := modelMeta.Field(i)
		dTag := newDynmgrmTag(tf.Tag)
		cn := getDBNameFromStructField(tf)
		isKey := false
		if dTag.pk {
			res.pk = dynmgrmKeyDefine{
				name:     cn,
				dataType: extractDBTypeFromStructField(tf),
			}
			isKey = true
		}
		if dTag.sk {
			res.sk = dynmgrmKeyDefine{
				name:     cn,
				dataType: extractDBTypeFromStructField(tf),
			}
			isKey = true
		}
		if !isKey {
			res.nonKey = append(res.nonKey, cn)
		}
		for _, ip := range dTag.indexProperty {
			switch ip.kind {
			case secondaryIndexKindGSI:
				sid, ok := res.gsi[ip.name]
				if !ok {
					sid = &dynmgrmSecondaryIndexDefine{}
					res.gsi[ip.name] = sid
				}
				if ip.pk {
					sid.pk = dynmgrmKeyDefine{
						name:     cn,
						dataType: extractDBTypeFromStructField(tf),
					}
				}
				if ip.sk {
					sid.sk = dynmgrmKeyDefine{
						name:     cn,
						dataType: extractDBTypeFromStructField(tf),
					}
				}
			case secondaryIndexKindLSI:
				res.lsi[ip.name] = &dynmgrmSecondaryIndexDefine{
					sk: dynmgrmKeyDefine{
						name:     cn,
						dataType: extractDBTypeFromStructField(tf),
					}}
			}
		}
		for _, np := range dTag.nonProjective {
			index, ok := res.lsi[np]
			if ok {
				index.nonProjectiveAttrs = append(index.nonProjectiveAttrs, cn)
				continue
			}
			index, ok = res.gsi[np]
			if ok {
				index.nonProjectiveAttrs = append(index.nonProjectiveAttrs, cn)
			}
		}
	}
	return res
}

func extractDBTypeFromStructField(field reflect.StructField) string {
	dbType := getDBTypeFromStructField(field)
	if dbType != "" {
		return dbType
	}
	switch field.Type.Kind() {
	case reflect.String:
		return KeySchemaDataTypeString.String()
	case reflect.Int, reflect.Float64:
		return KeySchemaDataTypeNumber.String()
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			return KeySchemaDataTypeBinary.String()
		}
	}
	return KeySchemaDataTypeString.String()
}
