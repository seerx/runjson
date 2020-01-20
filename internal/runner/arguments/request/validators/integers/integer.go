package integers

import "reflect"

type Integer interface {
	PrepareValue(val interface{}) error
	GreatThen(max int64, include bool) bool
	LessThen(min int64, include bool) bool
	Value() int64
}

var intMap = map[reflect.Kind]func() Integer{
	reflect.Int:    func() Integer { return &Int{} },
	reflect.Uint:   func() Integer { return &Uint{} },
	reflect.Int8:   func() Integer { return &Int8{} },
	reflect.Uint8:  func() Integer { return &Uint8{} },
	reflect.Int16:  func() Integer { return &Int16{} },
	reflect.Uint16: func() Integer { return &Uint16{} },
	reflect.Int32:  func() Integer { return &Int32{} },
	reflect.Uint32: func() Integer { return &Uint32{} },
	reflect.Int64:  func() Integer { return &Int64{} },
	reflect.Uint64: func() Integer { return &Uint64{} },
}

func FindMaker(kind reflect.Kind) func() Integer {
	if mk, ok := intMap[kind]; ok {
		return mk
	}
	return nil
}
