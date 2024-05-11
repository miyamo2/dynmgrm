package dynmgrm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

// ErrNestedStructHasIncompatibleAttributes occurs when the nested struct has incompatible attributes.
var ErrNestedStructHasIncompatibleAttributes = errors.New("nested struct has incompatible attributes")

// assignMapValueToReflectValue assigns the map type value to the reflect.Value
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
		err := assignInterfaceValueToReflectValue(tf.Type, vf, a)
		if err != nil {
			return err
		}
	}
	return nil
}

// assignInterfaceValueToReflectValue assigns the value to the reflect.Value
func assignInterfaceValueToReflectValue(rt reflect.Type, rv reflect.Value, value interface{}) error {
	if rv.CanAddr() {
		switch sc := rv.Addr().Interface().(type) {
		case sql.Scanner:
			return sc.Scan(value)
		}
	} else {
		switch sc := rv.Interface().(type) {
		case sql.Scanner:
			return sc.Scan(value)
		}
	}
	switch rt.Kind() {
	case reflect.String:
		str, ok := value.(string)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible string and %T", value))
		}
		rv.SetString(str)
	case reflect.Int:
		f64, ok := value.(float64)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible int and %T", value))
		}
		rv.Set(reflect.ValueOf(int(f64)))
	case reflect.Bool:
		b, ok := value.(bool)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible bool and %T", value))
		}
		rv.SetBool(b)
	case reflect.Float64:
		f64, ok := value.(float64)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible float64 and %T", value))
		}
		rv.SetFloat(f64)
	case reflect.Slice:
		if rt.Elem().Kind() != reflect.Uint8 {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible []byte and %T", value))
		}
		b, ok := value.([]byte)
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible []byte and %T", value))
		}
		rv.SetBytes(b)
	case reflect.Struct:
		mv, ok := value.(map[string]interface{})
		if !ok {
			return errors.Join(ErrNestedStructHasIncompatibleAttributes,
				fmt.Errorf("incompatible struct and %T", value))
		}
		err := assignMapValueToReflectValue(rt, rv, mv)
		if err != nil {
			return err
		}
	case reflect.Pointer:
		if value == nil {
			return nil
		}
		rv.Set(reflect.New(rt.Elem()))
		// NOTE: even return error, it will not be returned to the caller.
		// Only expect the attribute to be nil.
		assignInterfaceValueToReflectValue(rt.Elem(), rv.Elem(), value)
	}
	return nil
}
