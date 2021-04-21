package validators

import (
	"fmt"
)

//内置类型检查

// Validator 数据检查
type Validator interface {
	Check(val interface{}) error
}

// func ignoreCh(ch string) bool {
// 	return ch == " "
// }

func typeError(field, expected string) error {
	return fmt.Errorf("field %s expect type %s", field, expected)
}

func getFmt(field string, desc string,
	limitMax bool,
	max string,
	includeMax bool,
	limitMin bool,
	min string,
	includeMin bool, got string) string {

	errfmt := field + "'s " + desc

	minCond := ""
	if limitMin {
		if includeMin {
			minCond = "great then or equal " + min
		} else {
			minCond = "great then " + min
		}
	}

	maxCond := ""
	if limitMax {
		if includeMax {
			maxCond = "less then or equal " + max
		} else {
			maxCond = "less then " + max
		}
	}

	if minCond != "" {
		if maxCond != "" {
			errfmt = fmt.Sprintf("%s must %s and %s", errfmt, minCond, maxCond)
		} else {
			errfmt = fmt.Sprintf("%s must %s", errfmt, minCond)
		}
	} else {
		errfmt = fmt.Sprintf("%s must %s", errfmt, maxCond)
	}

	return errfmt + ", but got " + got
}
