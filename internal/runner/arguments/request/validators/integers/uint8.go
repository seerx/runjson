package integers

import (
	"errors"
	"reflect"
)

type Uint8 struct {
	val uint8
}

func (i *Uint8) Value() int64 {
	return int64(i.val)
}

func (i *Uint8) GreatThen(max int64, include bool) bool {
	val := uint8(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Uint8) LessThen(min int64, include bool) bool {
	val := uint8(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Uint8) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *uint8
		tmp, ok = val.(*uint8)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(uint8)
	}

	if !ok {
		return errors.New("uint8")
	}
	return nil
}
