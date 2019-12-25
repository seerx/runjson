package fieldmap

import "fmt"

type FieldMap map[string]int

func (fm *FieldMap) Requires(fields ...string) {
	for _, fd := range fields {
		if _, e := (*fm)[fd]; !e {
			panic(fmt.Errorf("%s is required", fd))
		}
	}
}

func (fm *FieldMap) Add(field string) {
	(*fm)[field] = 1
}

func (fm *FieldMap) HasFields(fields ...string) bool {
	for _, fd := range fields {
		if _, exists := (*fm)[fd]; !exists {
			return false
		}
	}
	return true
}
