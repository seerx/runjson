package validators

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/seerx/runjson/internal/runner/arguments/request/validators/integers"
)

// IntegerRange 检测整形范围
// 对应 tag 中的 limit 标签
// limit=0<$v  大于 0
// limit=$v<0  小于 0
// limit=-10<$v<0  大于 -10 小于 0
// 不允许出现 > 符号
// 大于小于 < 可以使用 <=  替换
type IntegerRange struct {
	field      string
	limitMax   bool
	max        int64
	includeMax bool

	limitMin   bool
	min        int64
	includeMin bool

	integerMaker func() integers.Integer

	errorFmt     string
	errorMessage string
}

func trim(val string) string {
	val = strings.TrimSpace(val)
	val = strings.Trim(val, "\t")
	val = strings.Trim(val, "\t")
	val = strings.Trim(val, "\n")
	return val
}

type rawRange struct {
	Min        string
	IncludeMin bool
	Max        string
	IncludeMax bool
}

func parseRange(exp string) (*rawRange, error) {
	ary := strings.Split(exp, ",")
	if len(ary) != 2 {
		return nil, fmt.Errorf("Invalid range expression:[%s]", exp)
	}
	rg := &rawRange{}
	min := trim(ary[0])
	if len(min) > 0 {
		rg.IncludeMin = min[:1] == "["
		if rg.IncludeMin || min[:1] == "(" {
			rg.Min = trim(min[1:])
		} else {
			rg.Min = min
		}
	}

	max := trim(ary[1])
	if len(max) > 0 {
		rg.IncludeMax = max[len(max)-1:] == "]"
		if rg.IncludeMax || max[len(max)-1:] == ")" {
			rg.Max = trim(max[:len(max)-1])
		} else {
			rg.Max = max
		}

	}

	return rg, nil
}

// CreateIntegerLimit 解析 limit 内容
// (1,2) [1,2] 1,2
func CreateIntegerLimit(fieldName, exp string, errorMessage string, mk func() integers.Integer, warnFn func(err error)) *IntegerRange {
	rg, err := parseRange(exp)
	if err != nil {
		warnFn(err)
		return nil
	}
	v := &IntegerRange{
		field:        fieldName,
		integerMaker: mk,
	}
	//v.
	if rg.Min != "" {
		intval, err := strconv.ParseInt(rg.Min, 0, 64)
		//intval, err := strconv.Atoi(rg.Min)
		if err != nil { // 发生错误
			warnFn(fmt.Errorf("Invalid range expression: [%s]: %w", exp, err))
		} else {
			v.limitMin = true
			v.min = intval
			v.includeMin = rg.IncludeMin
		}
	}
	if rg.Max != "" {
		intval, err := strconv.ParseInt(rg.Max, 0, 64)
		//intval, err := strconv.Atoi(rg.Max)
		if err != nil { // 发生错误
			warnFn(fmt.Errorf("Invalid range expression: [%s]: %w", exp, err))
		} else {
			v.limitMax = true
			v.max = intval
			v.includeMax = rg.IncludeMax
		}
	}
	v.errorFmt = getFmt(v.field, "value", v.limitMax, fmt.Sprintf("%d", v.max), v.includeMax,
		v.limitMin, fmt.Sprintf("%d", v.min), v.includeMin, "%v")
	v.errorMessage = errorMessage

	return v
}

func (v *IntegerRange) generateError(n interface{}) error {
	if v.errorMessage != "" {
		return errors.New(v.errorMessage)
	}
	return fmt.Errorf(v.errorFmt, n)
}

func (v *IntegerRange) Check(val interface{}) error {
	//var n int
	//var ok bool
	//if reflect.TypeOf(val).Kind() == reflect.Ptr {
	//	var tmp *int
	//	tmp, ok = val.(*int)
	//	if ok {
	//		n = *tmp
	//	}
	//} else {
	//	n, ok = val.(int)
	//}
	//
	//if !ok {
	//	return typeError(v.field, "int")
	//}

	it := v.integerMaker()
	if err := it.PrepareValue(val); err != nil {
		return typeError(v.field, err.Error())
	}

	if v.limitMax {
		// 限制了最大值
		if it.GreatThen(v.max, v.includeMax) {
			return v.generateError(it.Value())
		}
	}
	if v.limitMin {
		// 限制了最小值
		if it.LessThen(v.min, v.includeMin) {
			return v.generateError(it.Value())
		}
	}
	return nil
}
