package integers

import (
	"errors"
	"reflect"
)

type Int struct {
	val int
}

func (i *Int) Value() int64 {
	return int64(i.val)
}

func (i *Int) GreatThen(max int64, include bool) bool {
	val := int(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Int) LessThen(min int64, include bool) bool {
	val := int(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Int) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *int
		tmp, ok = val.(*int)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(int)
	}

	if !ok {
		return errors.New("int")
	}
	return nil
}
