package validators

import (
	"reflect"

	"github.com/seerx/runjson/internal/runner/arguments/request/validators/integers"

	"github.com/seerx/runjson/internal/reflects"
)

// GenerateValidators 生成验证信息
func GenerateValidators(typ reflect.Type, tag *reflects.ChainTag, warn func(err error)) []Validator {
	kind := typ.Kind()
	if typ.Kind() == reflect.Ptr {
		kind = typ.Elem().Kind()
	}
	validators := []Validator{}
	if tag.Limit != "" {
		created := false
		mk := integers.FindMaker(kind)
		if mk != nil {
			// 整形数据
			v := CreateIntegerLimit(tag.FieldName, tag.Limit, tag.Error, mk, warn)
			if v != nil {
				validators = append(validators, v)
				created = true
			}
		}
		if !created {
			if IsFloat(kind) {
				v := CreateFloatLimit(tag.FieldName, tag.Limit, tag.Error, warn)
				if v != nil {
					validators = append(validators, v)
					created = true
				}
			}
		}
		if !created {
			if IsString(kind) {
				v := CreateStringLimit(tag.FieldName, tag.Limit, tag.Error, warn)
				if v != nil {
					validators = append(validators, v)
					created = true
				}
			}
		}
	}
	if tag.Regexp != "" {
		if IsString(kind) {
			v := CreateRegexpValidator(tag.FieldName, tag.Limit, tag.Error)
			if v != nil {
				validators = append(validators, v)
			}
		}
	}
	return validators
}

// IsInt 是否整数类型
//func IsInt(kind reflect.Kind) bool {
//	return kind == reflect.Int ||
//		kind == reflect.Int8 ||
//		kind == reflect.Int16 ||
//		kind == reflect.Int32 ||
//		kind == reflect.Int64 ||
//		kind == reflect.Uint ||
//		kind == reflect.Uint8 ||
//		kind == reflect.Uint16 ||
//		kind == reflect.Uint32 ||
//		kind == reflect.Uint64
//}

// IsFloat 是否浮点类型
func IsFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 ||
		kind == reflect.Float64
}

// IsString 是否字符串类型
func IsString(kind reflect.Kind) bool {
	return kind == reflect.String
}
