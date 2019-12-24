package reflects

import (
	"reflect"
	"strings"
)

const (
	TAG     = "c"
	tagDesc = "desc"
	//tagPrefix       = "prefix"
	tagLimit        = "limit"
	tagRegexp       = "regexp"
	tagErrorMessage = "error"
	tagRequire      = "require"
	tagDeprecated   = "deprecated"
	//tagExplodeParams = "explode"

)

// ChainTag
type ChainTag struct {
	FieldName   string // 字段名称,使用 json 定义，如果没有则使用 fieldName
	Require     bool   // 必填字段
	Description string // 说明
	Limit       string // 限制
	Regexp      string // 正则表达式限制
	Error       string // 设置的错误提示
	Deprecated  bool   // 建议不要使用
}

// ParseTag 解析 Tag，如果  json tag 设置为 -,则返回 nil,即忽略此字段
func ParseTag(field *reflect.StructField) *ChainTag {
	gTag := &ChainTag{}
	tag := field.Tag

	gTag.FieldName = parseFieldName(&tag, field)
	if gTag.FieldName == "" {
		return nil
	}
	gqlTag := tag.Get(TAG)
	if gqlTag != "" {
		mp := map[string]string{}
		ary := strings.Split(gqlTag, ",")
		for _, item := range ary {
			p := strings.Index(item, ":")
			if p < 0 {
				mp[item] = item
			} else {
				mp[item[:p]] = item[p+1:]
			}
			//sub := strings.Split(item, "=")
			//if len(sub) == 1 {
			//	mp[sub[0]] = sub[0]
			//}
			//if len(sub) == 2 {
			//	mp[sub[0]] = sub[1]
			//}
		}

		//gTag.Prefix = mp[tagPrefix] //  parseGqlPrefix(mp)
		gTag.Description = mp[tagDesc]
		gTag.Regexp = mp[tagRegexp]
		gTag.Limit = mp[tagLimit]
		gTag.Error = mp[tagErrorMessage]
		_, gTag.Require = mp[tagRequire]
		_, gTag.Deprecated = mp[tagDeprecated]
	}

	return gTag
}

// parseGqlPrefix 解析函数前缀
//func parseGqlPrefix(pairs map[string]string) string {
//	return pairs[tagPrefix]
//}

// parseFieldName 解析字段名称
func parseFieldName(tag *reflect.StructTag, field *reflect.StructField) string {
	name := tag.Get("json")
	if name == "" {
		name = field.Name
	}
	if name == "-" {
		// json 中的忽略
		return ""
	}
	ary := strings.Split(name, ",")
	if len(ary) == 1 {
		return name
	}

	for _, item := range ary {
		if item != "omitempty" {
			return item
		}
	}

	return ""
}
