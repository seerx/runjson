package integers

import (
	"errors"
	"reflect"
)

type Uint16 struct {
	val uint16
}

func (i *Uint16) Value() int64 {
	return int64(i.val)
}

func (i *Uint16) GreatThen(max int64, include bool) bool {
	val := uint16(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Uint16) LessThen(min int64, include bool) bool {
	val := uint16(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Uint16) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *uint16
		tmp, ok = val.(*uint16)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(uint16)
	}

	if !ok {
		return errors.New("uint16")
	}
	return nil
}
