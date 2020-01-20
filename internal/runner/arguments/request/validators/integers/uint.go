package integers

import (
	"errors"
	"reflect"
)

type Uint struct {
	val uint
}

func (i *Uint) Value() int64 {
	return int64(i.val)
}

func (i *Uint) GreatThen(max int64, include bool) bool {
	val := uint(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Uint) LessThen(min int64, include bool) bool {
	val := uint(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Uint) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *uint
		tmp, ok = val.(*uint)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(uint)
	}

	if !ok {
		return errors.New("uint")
	}
	return nil
}
