package integers

import (
	"errors"
	"reflect"
)

type Int64 struct {
	val int64
}

func (i *Int64) Value() int64 {
	return int64(i.val)
}

func (i *Int64) GreatThen(max int64, include bool) bool {
	val := int64(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Int64) LessThen(min int64, include bool) bool {
	val := int64(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Int64) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *int64
		tmp, ok = val.(*int64)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(int64)
	}

	if !ok {
		return errors.New("Int64")
	}
	return nil
}
