package validators

import (
	"reflect"

	"github.com/seerx/runjson/internal/reflects"
)

// GenerateValidators 生成验证信息
func GenerateValidators(typ reflect.Type, tag *reflects.ChainTag) []Validator {
	validators := []Validator{}
	if tag.Limit != "" {
		if IsInt(typ.Kind()) {
			v := CreateIntegerLimit(tag.FieldName, tag.Limit, tag.Error)
			if v != nil {
				validators = append(validators, v)
			}
		} else if IsFloat(typ.Kind()) {
			v := CreateFloatLimit(tag.FieldName, tag.Limit, tag.Error)
			if v != nil {
				validators = append(validators, v)
			}
		} else if IsString(typ.Kind()) {
			v := CreateStringLimit(tag.FieldName, tag.Limit, tag.Error)
			if v != nil {
				validators = append(validators, v)
			}
		}
	}
	if tag.Regexp != "" {
		if IsString(typ.Kind()) {
			v := CreateRegexpValidator(tag.FieldName, tag.Limit, tag.Error)
			if v != nil {
				validators = append(validators, v)
			}
		}
	}
	return validators
}

// IsInt 是否整数类型
func IsInt(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64
}

// IsFloat 是否浮点类型
func IsFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 ||
		kind == reflect.Float64
}

// IsString 是否字符串类型
func IsString(kind reflect.Kind) bool {
	return kind == reflect.String
}
