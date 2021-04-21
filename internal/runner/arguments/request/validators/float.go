package validators

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// FloatRange 检测整形范围
// 对应 tag 中的 limit 标签
// limit=0<$v  大于 0
// limit=$v<0  小于 0
// limit=-10<$v<0  大于 -10 小于 0
// 不允许出现 > 符号
// 大于小于 < 可以使用 <=  替换
type FloatRange struct {
	field      string
	limitMax   bool
	max        float64
	includeMax bool

	limitMin   bool
	min        float64
	includeMin bool

	errorFmt     string
	errorMessage string
}

// CreateFloatLimit 解析 limit 内容
func CreateFloatLimit(fieldName string, exp string, errorMessage string, warnFn func(err error)) *FloatRange {
	rg, err := parseRange(exp)
	if err != nil {
		warnFn(err)
		return nil
	}
	v := &FloatRange{
		field: fieldName,
	}

	if rg.Min != "" {
		val, err := strconv.ParseFloat(rg.Min, 64)
		//intval, err := strconv.Atoi(rg.Min)
		if err != nil { // 发生错误
			//warnFn(err)
			warnFn(fmt.Errorf("invalid range expression: [%s]: %w", exp, err))
		} else {
			v.limitMin = true
			v.min = val
			v.includeMin = rg.IncludeMin
		}
	}
	if rg.Max != "" {
		val, err := strconv.ParseFloat(rg.Max, 64)
		//intval, err := strconv.Atoi(rg.Max)
		if err != nil { // 发生错误
			//warnFn(err)
			warnFn(fmt.Errorf("invalid range expression: [%s]: %w", exp, err))
		} else {
			v.limitMax = true
			v.max = val
			v.includeMax = rg.IncludeMax
		}
	}

	//vp := strings.Index(exp, "$v")
	//if vp < 0 {
	//	// 没有找到 $v
	//	return nil
	//}
	//v := &FloatRange{
	//	field: fieldName,
	//}
	//
	//// 计算 vp 前面有字符
	//for n := vp - 1; n >= 0; n-- {
	//	ch := exp[n : n+1]
	//	if ignoreCh(ch) {
	//		// 忽略空格
	//		continue
	//	}
	//	if v.limitMin {
	//		// 解析数字
	//		val := exp[:n+1]
	//		floatVal, err := strconv.ParseFloat(val, 64)
	//		if err != nil {
	//			panic(fmt.Errorf("%s' tag %s cann't convert to float %s", v.field, val, err.Error()))
	//		}
	//		v.min = floatVal
	//		break
	//	}
	//	if ch == "=" {
	//		// 出现等号
	//		v.includeMin = true
	//		continue
	//	}
	//	if ch == "<" {
	//		// 前面出现小于
	//		v.limitMin = true
	//	}
	//}
	//
	//for n := vp + 2; n < len(exp); n++ {
	//	ch := exp[n : n+1]
	//	if ignoreCh(ch) {
	//		// 忽略空格
	//		continue
	//	}
	//	if ch == "=" {
	//		// 出现等号
	//		v.includeMax = true
	//		continue
	//	}
	//
	//	if v.limitMax {
	//		// 解析数字
	//		val := exp[n:]
	//		floatVal, err := strconv.ParseFloat(val, 64)
	//		if err != nil {
	//			panic(fmt.Errorf("%s' tag %s cann't convert to integer %s", v.field, val, err.Error()))
	//		}
	//		v.max = floatVal
	//		break
	//	}
	//
	//	if ch == "<" {
	//		// 前面出现小于
	//		v.limitMax = true
	//	}
	//}
	v.errorFmt = getFmt(v.field, "value", v.limitMax, fmt.Sprintf("%f", v.max), v.includeMax,
		v.limitMin, fmt.Sprintf("%f", v.min), v.includeMin, "%f")
	v.errorMessage = errorMessage
	return v
}

func (v *FloatRange) generateError(n float64) error {
	if v.errorMessage != "" {
		return errors.New(v.errorMessage)
	}
	return fmt.Errorf(v.errorFmt, n)
}

// Check 检查数据
func (v *FloatRange) Check(val interface{}) error {
	var n float64
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *float64
		tmp, ok = val.(*float64)
		if ok {
			n = *tmp
		}
	} else {
		n, ok = val.(float64)
	}
	if !ok {
		return typeError(v.field, "float64")
	}
	// }

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
