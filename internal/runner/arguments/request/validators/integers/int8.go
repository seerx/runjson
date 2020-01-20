package integers

import (
	"errors"
	"reflect"
)

type Int8 struct {
	val int8
}

func (i *Int8) Value() int64 {
	return int64(i.val)
}

func (i *Int8) GreatThen(max int64, include bool) bool {
	val := int8(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Int8) LessThen(min int64, include bool) bool {
	val := int8(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Int8) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *int8
		tmp, ok = val.(*int8)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(int8)
	}

	if !ok {
		return errors.New("Int8")
	}
	return nil
}
