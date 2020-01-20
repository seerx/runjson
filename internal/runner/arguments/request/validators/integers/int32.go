package integers

import (
	"errors"
	"reflect"
)

type Int32 struct {
	val int32
}

func (i *Int32) Value() int64 {
	return int64(i.val)
}

func (i *Int32) GreatThen(max int64, include bool) bool {
	val := int32(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Int32) LessThen(min int64, include bool) bool {
	val := int32(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Int32) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *int32
		tmp, ok = val.(*int32)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(int32)
	}

	if !ok {
		return errors.New("Int32")
	}
	return nil
}
