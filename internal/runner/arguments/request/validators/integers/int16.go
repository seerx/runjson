package integers

import (
	"errors"
	"reflect"
)

type Int16 struct {
	val int16
}

func (i *Int16) Value() int64 {
	return int64(i.val)
}

func (i *Int16) GreatThen(max int64, include bool) bool {
	val := int16(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Int16) LessThen(min int64, include bool) bool {
	val := int16(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Int16) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *int16
		tmp, ok = val.(*int16)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(int16)
	}

	if !ok {
		return errors.New("Int16")
	}
	return nil
}
