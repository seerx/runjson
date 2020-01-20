package integers

import (
	"errors"
	"reflect"
)

type Uint64 struct {
	val uint64
}

func (i *Uint64) Value() int64 {
	return int64(i.val)
}

func (i *Uint64) GreatThen(max int64, include bool) bool {
	val := uint64(max)
	if include {
		return i.val > val
	}
	return i.val >= val
}

func (i *Uint64) LessThen(min int64, include bool) bool {
	val := uint64(min)
	if include {
		return i.val < val
	}
	return i.val <= val
}

func (i *Uint64) PrepareValue(val interface{}) error {
	//var n int
	var ok bool
	if reflect.TypeOf(val).Kind() == reflect.Ptr {
		var tmp *uint64
		tmp, ok = val.(*uint64)
		if ok {
			i.val = *tmp
		}
	} else {
		i.val, ok = val.(uint64)
	}

	if !ok {
		return errors.New("uint64")
	}
	return nil
}
