package input

import "reflect"

type Requirement struct {
}

var typeOfRequirement = reflect.TypeOf(Requirement{})

func IsRequirementArgument(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		return typ.Elem() == typeOfRequirement
	}
	return typ == typeOfRequirement
}
