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
	PK   bool
	SK   bool
	Name string
	Kind secondaryIndexKind
}

type dynmgrmTag struct {
	PK            bool
	SK            bool
	IndexProperty []secondaryIndexProperty
	NonProjective []string
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
			if res.SK {
				continue
			}
			res.PK = true
		case "sk":
			if res.PK {
				continue
			}
			res.SK = true
		case "gsi-pk":
			iprp := secondaryIndexProperty{
				PK:   true,
				Name: kv[1],
				Kind: secondaryIndexKindGSI,
			}
			res.IndexProperty = append(res.IndexProperty, iprp)
		case "gsi-sk":
			iprp := secondaryIndexProperty{
				SK:   true,
				Name: kv[1],
				Kind: secondaryIndexKindGSI,
			}
			res.IndexProperty = append(res.IndexProperty, iprp)
		case "lsi-sk":
			iprp := secondaryIndexProperty{
				Name: kv[1],
				SK:   true,
				Kind: secondaryIndexKindLSI,
			}
			res.IndexProperty = append(res.IndexProperty, iprp)
		case "non-projective":
			npl := strings.ReplaceAll(strings.ReplaceAll(kv[1], "[", ""), "]", "")
			for _, np := range strings.Split(npl, ",") {
				res.NonProjective = append(res.NonProjective, np)
			}
		}
	}
	return res
}

type dynmgrmKeyDefine struct {
	Name     string
	DataType string
}

type dynmgrmSecondaryIndexDefine struct {
	PK                 dynmgrmKeyDefine
	SK                 dynmgrmKeyDefine
	NonProjectiveAttrs []string
}

type dynmgrmTableDefine struct {
	PK         dynmgrmKeyDefine
	SK         dynmgrmKeyDefine
	NonKeyAttr []string
	GSI        map[string]*dynmgrmSecondaryIndexDefine
	LSI        map[string]*dynmgrmSecondaryIndexDefine
}

func newDynmgrmTableDefine(modelMeta reflect.Type) dynmgrmTableDefine {
	res := dynmgrmTableDefine{
		NonKeyAttr: make([]string, 0, modelMeta.NumField()),
		GSI:        make(map[string]*dynmgrmSecondaryIndexDefine),
		LSI:        make(map[string]*dynmgrmSecondaryIndexDefine),
	}

	nonProjectiveAttrsMap := make(map[string]*[]string)
	for i := 0; i < modelMeta.NumField(); i++ {
		tf := modelMeta.Field(i)
		dTag := newDynmgrmTag(tf.Tag)
		cn := getDBNameFromStructField(tf)
		isKey := false
		if dTag.PK {
			res.PK = dynmgrmKeyDefine{
				Name:     cn,
				DataType: extractDBTypeFromStructField(tf),
			}
			isKey = true
		}
		if dTag.SK {
			res.SK = dynmgrmKeyDefine{
				Name:     cn,
				DataType: extractDBTypeFromStructField(tf),
			}
			isKey = true
		}
		if !isKey {
			res.NonKeyAttr = append(res.NonKeyAttr, cn)
		}
		for _, ip := range dTag.IndexProperty {
			switch ip.Kind {
			case secondaryIndexKindGSI:
				sid, ok := res.GSI[ip.Name]
				if !ok {
					sid = &dynmgrmSecondaryIndexDefine{}
					res.GSI[ip.Name] = sid
				}
				if ip.PK {
					sid.PK = dynmgrmKeyDefine{
						Name:     cn,
						DataType: extractDBTypeFromStructField(tf),
					}
				}
				if ip.SK {
					sid.SK = dynmgrmKeyDefine{
						Name:     cn,
						DataType: extractDBTypeFromStructField(tf),
					}
				}
			case secondaryIndexKindLSI:
				sid, ok := res.LSI[ip.Name]
				if !ok {
					sid = &dynmgrmSecondaryIndexDefine{}
					res.LSI[ip.Name] = sid
				}
				sid.SK = dynmgrmKeyDefine{
					Name:     cn,
					DataType: extractDBTypeFromStructField(tf),
				}
			}
		}
		for _, np := range dTag.NonProjective {
			list, ok := nonProjectiveAttrsMap[np]
			if !ok {
				physicalList := make([]string, 0)
				list = &physicalList
				nonProjectiveAttrsMap[np] = list
			}
			*list = append(*list, cn)
		}
	}
	for idxn, list := range nonProjectiveAttrsMap {
		index, ok := res.LSI[idxn]
		if !ok {
			index, ok = res.GSI[idxn]
			if !ok {
				continue
			}
		}
		// list will never be nil
		if pl := *list; len(pl) > 0 {
			index.NonProjectiveAttrs = append(index.NonProjectiveAttrs, pl...)
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
