package validators

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// StringLengthRange 检测整形范围
// 对应 tag 中的 limit 标签
// limit=0<$v  大于 0
// limit=$v<0  小于 0
// limit=-10<$v<0  大于 -10 小于 0
// 不允许出现 > 符号
// 大于小于 < 可以使用 <=  替换
type StringLengthRange struct {
	field      string
	limitMax   bool
	max        int
	includeMax bool

	limitMin   bool
	min        int
	includeMin bool

	errorFmt     string
	errorMessage string
}

// CreateStringLimit 解析 limit 内容
func CreateStringLimit(fieldName string, exp string, errorMessage string, warnFn func(err error)) *StringLengthRange {
	rg, err := parseRange(exp)
	if err != nil {
		warnFn(err)
		return nil
	}

	v := &StringLengthRange{
		field: fieldName,
	}

	if rg.Min != "" {
		//intval, err := strconv.ParseInt(rg.Min, 0, 64)
		intval, err := strconv.Atoi(rg.Min)
		if err != nil { // 发生错误
			//warnFn(err)
			warnFn(fmt.Errorf("Invalid range expression: [%s]: %w", exp, err))
		} else {
			v.limitMin = true
			v.min = intval
			v.includeMin = rg.IncludeMin
		}
	}
	if rg.Max != "" {
		//intval, err := strconv.ParseInt(rg.Max, 0, 64)
		intval, err := strconv.Atoi(rg.Max)
		if err != nil { // 发生错误
			//warnFn(err)
			warnFn(fmt.Errorf("Invalid range expression: [%s]: %w", exp, err))
		} else {
			v.limitMax = true
			v.max = intval
			v.includeMax = rg.IncludeMax
		}
	}

	v.errorFmt = getFmt(v.field, "length", v.limitMax, fmt.Sprintf("%d", v.max), v.includeMax,
		v.limitMin, fmt.Sprintf("%d", v.min), v.includeMin, "%d")
	v.errorMessage = errorMessage
	return v
}

func (v *StringLengthRange) generateError(n int) error {
	if v.errorMessage != "" {
		return errors.New(v.errorMessage)
	}
	return fmt.Errorf(v.errorFmt, n)
}

func (v *StringLengthRange) Check(val interface{}) error {
	var str string
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var s *string
		s, ok = val.(*string)
		if ok {
			str = *s
		}
	} else {
		str, ok = val.(string)
	}
	if !ok {
		return typeError(v.field, "string")
	}
	n := len(str)
	if v.limitMax {
		// 限制了最大值
		if v.includeMax {
			if n > v.max {
				return v.generateError(n)
			}
		} else {
			if n >= v.max {
				return v.generateError(n)
			}
		}
	}
	if v.limitMin {
		// 限制了最小值
		if v.includeMin {
			if n < v.min {
				return v.generateError(n)
			}
		} else {
			if n <= v.min {
				return v.generateError(n)
			}
		}
	}
	return nil
}
