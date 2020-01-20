package integers

import (
	"errors"
	"reflect"
)

type Uint32 struct {
	val uint32
}

func (i *Uint32) Value() int64 {
	return int64(i.val)
}

func (i *Uint32) GreatThen(max int64, include bool) bool {
	val := uint32(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Uint32) LessThen(min int64, include bool) bool {
	val := uint32(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Uint32) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *uint32
		tmp, ok = val.(*uint32)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(uint32)
	}

	if !ok {
		return errors.New("uint32")
	}
	return nil
}
