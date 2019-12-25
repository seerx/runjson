package intf

type FieldMap map[string]int

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
