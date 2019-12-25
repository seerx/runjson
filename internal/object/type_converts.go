package object

import (
	"fmt"
	"reflect"
)

// typ 使用 Refernce
func tryToConvert(fieldName string, typ reflect.Type, val interface{}) (reflect.Value, error) {
	//vType := reflect.TypeOf(val)
	//if vType.Kind() == reflect.Ptr {
	//	vType = vType.Elem()
	//}

	var nilVal reflect.Value

	//if vType.ConvertibleTo(typ) {
	//var (
	//	Int   int
	//	Float float64
	//)
	switch typ.Kind() {
	case reflect.Int:
		nilVal = reflect.ValueOf(new(int))
		Int, ok := val.(int)
		if ok {
			if res := int2Type(typ, Int); !res.IsNil() {
				return res, nil
			}
		} else {
			Float, ok := val.(float64)
			if ok {
				if res := float2Type(typ, Float); !res.IsNil() {
					return res, nil
				}
			}
		}
	case reflect.Float64, reflect.Float32:
		nilVal = reflect.ValueOf(new(float64))
		Float, ok := val.(float64)
		//return reflect.ValueOf(val).Convert(typ)
		if ok {
			if res := float2Type(typ, Float); !res.IsNil() {
				return res, nil
			}
		}
	case reflect.String:
		nilVal = reflect.ValueOf(new(string))
		str, ok := val.(string)
		if ok {
			return reflect.ValueOf(&str), nil
		}
	case reflect.Bool:
		nilVal = reflect.ValueOf(new(bool))
		b, ok := val.(bool)
		if ok {
			return reflect.ValueOf(&b), nil
		}
	}
	//}
	// 数据类型不匹配
	return nilVal, fmt.Errorf("[%s] Invalid data type, expect %s", fieldName, nilVal.Elem().Type().Name())
	// string bool 直接返回
	//return reflect.ValueOf(val).Convert(typ)
}

func float2Type(typ reflect.Type, val float64) reflect.Value {
	// 可以转换
	switch typ.Kind() {
	case reflect.Int:
		n := int(val)
		return reflect.ValueOf(&n)
	case reflect.Int8:
		n := int8(val)
		return reflect.ValueOf(&n)
	case reflect.Int16:
		n := int16(val)
		return reflect.ValueOf(&n)
	case reflect.Int32:
		n := int32(val)
		return reflect.ValueOf(&n)
	case reflect.Int64:
		n := int64(val)
		return reflect.ValueOf(&n)

	case reflect.Uint:
		return reflect.ValueOf(&val)
	case reflect.Uint8:
		n := uint8(val)
		return reflect.ValueOf(&n)
	case reflect.Uint16:
		n := uint16(val)
		return reflect.ValueOf(&n)
	case reflect.Uint32:
		n := uint32(val)
		return reflect.ValueOf(&n)
	case reflect.Uint64:
		n := uint64(val)
		return reflect.ValueOf(&n)

	case reflect.Float32:
		f := float32(val)
		return reflect.ValueOf(&f)
	case reflect.Float64:
		return reflect.ValueOf(&val)
	}
	return reflect.ValueOf(nil)
}

func int2Type(typ reflect.Type, val int) reflect.Value {
	// 可以转换
	switch typ.Kind() {
	case reflect.Int:
		return reflect.ValueOf(&val)
	case reflect.Int8:
		n := int8(val)
		return reflect.ValueOf(&n)
	case reflect.Int16:
		n := int16(val)
		return reflect.ValueOf(&n)
	case reflect.Int32:
		n := int32(val)
		return reflect.ValueOf(&n)
	case reflect.Int64:
		n := int64(val)
		return reflect.ValueOf(&n)

	case reflect.Uint:
		return reflect.ValueOf(&val)
	case reflect.Uint8:
		n := uint8(val)
		return reflect.ValueOf(&n)
	case reflect.Uint16:
		n := uint16(val)
		return reflect.ValueOf(&n)
	case reflect.Uint32:
		n := uint32(val)
		return reflect.ValueOf(&n)
	case reflect.Uint64:
		n := uint64(val)
		return reflect.ValueOf(&n)

	case reflect.Float32:
		f := float32(val)
		return reflect.ValueOf(&f)
	case reflect.Float64:
		f := float64(val)
		return reflect.ValueOf(&f)
	}
	return reflect.ValueOf(nil)
}

//func int2Type(typ reflect.Type, val int) interface{} {
//	switch typ.Kind() {
//	case reflect.Int:
//		return int(val)
//	case reflect.Int8:
//		return int8(val)
//	case reflect.Int32:
//		return int32(val)
//	case reflect.Int64:
//		return int64(val)
//
//	case reflect.Uint:
//		return uint(val)
//	case reflect.Uint8:
//		return uint8(val)
//	case reflect.Uint16:
//		return uint16(val)
//	case reflect.Uint32:
//		return uint32(val)
//	case reflect.Uint64:
//		return uint64(val)
//
//	case reflect.Float32:
//		return float32(val)
//	case reflect.Float64:
//		return val
//	}
//
//	return nil
//}
