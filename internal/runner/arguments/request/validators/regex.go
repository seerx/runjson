package validators

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// RegexValidator 正则表达式验证
// 对应 tag 中的 regexp
type RegexpValidator struct {
	field        string
	regstr       string
	errorMessage string
	regex        *regexp.Regexp
}

func (v *RegexpValidator) Check(val interface{}) error {
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
		return fmt.Errorf("%s's input value is invalid", v.field)
		//return typeError("string")
	}
	found := v.regex.FindString(str)
	//v.regex.MatchString(str)
	if found != str {
		//if !v.regex.Match([]byte(str)) {
		if v.errorMessage != "" {
			return errors.New(v.errorMessage)
		}
		return fmt.Errorf("%s's value do not match with regular expression %s", v.field, v.regstr)
	}
	return nil
}

func CreateRegexpValidator(fieldName string, exp string, errMessage string) *RegexpValidator {
	reg, err := regexp.Compile(exp)
	if err != nil {
		panic(fmt.Errorf("Invalid regular expression %s = %s", fieldName, exp))
		return nil
	}
	return &RegexpValidator{
		regstr:       exp,
		field:        fieldName,
		regex:        reg,
		errorMessage: errMessage,
	}
}
