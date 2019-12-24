package reflects

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type Location struct {
	Package string
	Func    string
	Struct  string
}

func (l *Location) Equals(o *Location) bool {
	if o != nil {
		return l.String() == o.String()
	}
	return false
}

func (l *Location) String() string {
	if l.Struct == "" {
		if l.Func == "" {
			return l.Package
		}
		return fmt.Sprintf("%s func:%s", l.Package, l.Func)
	}

	if l.Func == "" {
		return fmt.Sprintf("%s struct:%s", l.Package, l.Struct)
	}
	return fmt.Sprintf("%s %s.%s", l.Package, l.Struct, l.Func)
}

// ParseFuncLocation 解析函数信息
func ParseFuncLocation(aFunc interface{}) *Location {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(aFunc).Pointer()).Name()

	p := strings.LastIndex(fn, ".")

	if p > 0 {
		return &Location{
			Func:    fn[p+1:],
			Package: fn[0:p],
		}
	}

	return nil
}

func ParseStructFuncLocation(structType reflect.Type, method reflect.Method) *Location {
	st := ParseType(&Location{}, structType)
	return &Location{
		Package: st.Package,
		Func:    method.Name,
		Struct:  st.Name,
	}
}

func ParseStructLocation(structType reflect.Type) *Location {
	st := ParseType(&Location{}, structType)
	return &Location{
		Package: st.Package,
		Func:    "",
		Struct:  st.Name,
	}
}
