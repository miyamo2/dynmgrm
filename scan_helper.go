package dynmgrm

import (
	"database/sql"
	"fmt"
	"reflect"
)

func assignMapValueToReflectValue(rt reflect.Type, rv reflect.Value, mv map[string]interface{}) error {
	for i := 0; i < rt.NumField(); i++ {
		tf := rt.Field(i)
		vf := func() reflect.Value {
			if rv.Kind() == reflect.Pointer {
				return rv.Elem().Field(i)
			}
			return rv.Field(i)
		}()
		name := getDBNameFromStructField(tf)
		a, ok := mv[name]
		if !ok {
			continue
		}
		assignInterfaceValueToReflectValue(tf.Type, vf, a)
	}
	return nil
}

func assignInterfaceValueToReflectValue(rt reflect.Type, rv reflect.Value, value interface{}) error {
	if rv.CanAddr() {
		switch ptr := rv.Addr().Interface().(type) {
		case sql.Scanner:
			if err := ptr.Scan(value); err != nil {
				return err
			}
		}
	}
	switch rt.Kind() {
	case reflect.String:
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("incompatible string and %T", value)
		}
		rv.SetString(str)
	case reflect.Int:
		f64, ok := value.(float64)
		if !ok {
			return fmt.Errorf("incompatible int and %T", value)
		}
		rv.Set(reflect.ValueOf(int(f64)))
	case reflect.Bool:
		b, ok := value.(bool)
		if !ok {
			return fmt.Errorf("incompatible bool and %T", value)
		}
		rv.SetBool(b)
	case reflect.Float64:
		f64, ok := value.(float64)
		if !ok {
			return fmt.Errorf("incompatible float64 and %T", value)
		}
		rv.SetFloat(f64)
	case reflect.Slice:
		if rt.Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("incompatible []byte and %T", value)
		}
		b, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("incompatible []byte and %T", value)
		}
		rv.SetBytes(b)
	case reflect.Struct:
		mv, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("incompatible struct and %T", value)
		}
		assignMapValueToReflectValue(rt, rv, mv)
	case reflect.Pointer:
		if value == nil {
			return nil
		}
		rv.Set(reflect.New(rt.Elem()))
		if err := assignInterfaceValueToReflectValue(rt.Elem(), rv.Elem(), value); err != nil {
			return err
		}
	}
	return nil
}
