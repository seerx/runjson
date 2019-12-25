package reflects

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/seerx/runjson/internal/util"
)

// TypeInfo 解析 reflect.Type 的信息
type TypeInfo struct {
	Raw       reflect.Type // 原始类型
	Reference reflect.Type // 指向类型，指针、数组、Map 等指向的类型

	// 可以看到的类型名称
	Name    string // 类型名称
	Package string // 所属包

	IsRawPtr   bool // 原始类型是否指针
	IsRawSlice bool // 原始类型是都是数组

	IsSliceItemIsPtr bool // 数组类型是否指针
	IsPrimitive      bool // 指向类型 是否原生类型
	IsStruct         bool // 指向类型 是否结构体

	IsInterface bool // 是否接口
	IsFunc      bool // 是否函数
}

func (ti *TypeInfo) Path() string {
	pkg := strings.ReplaceAll(ti.Package, ".", "_")
	pkg = strings.ReplaceAll(pkg, "/", "__")
	return fmt.Sprintf("%s_%s", ti.Name, pkg)
}

func (ti *TypeInfo) ID() string {
	return util.MD5(ti.Package, ti.Name)
}

// TypeName
func (ti *TypeInfo) TypeName() string {
	if ti.IsPrimitive {
		kd := ti.Reference.Kind()
		if kd == reflect.Bool {
			return "bool"
		}
		if kd == reflect.Int ||
			kd == reflect.Int8 ||
			kd == reflect.Int16 ||
			kd == reflect.Int32 ||
			kd == reflect.Int64 ||
			kd == reflect.Uint ||
			kd == reflect.Uint8 ||
			kd == reflect.Uint16 ||
			kd == reflect.Uint32 ||
			kd == reflect.Uint64 {
			return "int"
		}
		if kd == reflect.Float32 || kd == reflect.Float64 {
			return "float"
		}
		if kd == reflect.String {
			return "string"
		}
	}
	return ti.Reference.Name()
}

// FieldInfo 解析 struct 的 field 的 reflect.Type 信息
type FieldInfo struct {
	TypeInfo
	Field *reflect.StructField
}

// ParseField 解析结构字段类型
func ParseField(loc *Location, field *reflect.StructField) *FieldInfo {
	ti := ParseType(loc, field.Type)
	return &FieldInfo{
		TypeInfo: *ti,
		Field:    field,
	}
}

var (
	supportsPrimitive = map[reflect.Kind]int{
		reflect.Bool: 0,

		reflect.Int:    1,
		reflect.Int8:   1,
		reflect.Int16:  1,
		reflect.Int32:  1,
		reflect.Int64:  1,
		reflect.Uint:   1,
		reflect.Uint8:  1,
		reflect.Uint16: 1,
		reflect.Uint32: 1,
		reflect.Uint64: 1,

		reflect.Float32: 2,
		reflect.Float64: 2,

		reflect.String: 3,
	}
	supportsAdvance = map[reflect.Kind]int{
		reflect.Struct: 4,

		reflect.Ptr:   5,
		reflect.Slice: 6,
		//}
		//
		//supportsInterface = map[reflect.Kind]int{
		reflect.Interface: 7,
	}
)

var supports = map[reflect.Kind]int{}

func checkSupport(loc *Location, name string, typ reflect.Type) {
	if _, s := supportsPrimitive[typ.Kind()]; !s {
		// b不是原生类型
		if _, s := supportsAdvance[typ.Kind()]; !s {
			// 不是支持的类型
			panic(fmt.Errorf("[%s] is not support locate %s", name, loc.String()))
		}
	}
}

func checkPtrSupport(loc *Location, name string, typ reflect.Type) {
	if _, s := supportsPrimitive[typ.Kind()]; !s {
		if typ.Kind() != reflect.Struct {
			// 不是基础类型且不是结构 panic
			panic(fmt.Errorf("[%s] is not support locate %s", name, loc.String()))
		}
	}
}

// ParseType 解析类型
func ParseType(loc *Location, typ reflect.Type) *TypeInfo {
	//typName := typ.Name()
	//checkSupport(loc, typName, typ)

	kd := typ.Kind()
	info := &TypeInfo{
		Raw:       typ,
		Reference: typ,
		//Name:             typName,
		//Package:          typ.PkgPath(),
		IsRawPtr:         kd == reflect.Ptr,
		IsRawSlice:       kd == reflect.Slice,
		IsSliceItemIsPtr: false,

		//IsPrimitive: typ.PkgPath() == "",
		//IsStruct:    false,
		IsInterface: kd == reflect.Interface,
		IsFunc:      kd == reflect.Func,
	}

	//var ref = typ
	if info.IsRawPtr {
		// 指针
		info.Reference = typ.Elem()
		//checkPtrSupport(loc, info.Reference.Name(), info.Reference)
	} else if typ.Kind() == reflect.Slice {
		// 数组(切片)
		info.Reference = typ.Elem()
		if info.Reference.Kind() == reflect.Ptr {
			// 数组元素是指针
			info.IsSliceItemIsPtr = true
			info.Reference = info.Reference.Elem()
		}

	}

	info.Name = info.Reference.Name()
	info.Package = info.Reference.PkgPath()
	info.IsStruct = info.Reference.Kind() == reflect.Struct
	info.IsPrimitive = info.Package == ""

	checkPtrSupport(loc, info.Name, info.Reference)

	return info
}
