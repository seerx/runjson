package reflects

import (
	"reflect"
	"strings"
)

const (
	TAG     = "rj"
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

func isValidateToken(token string) bool {
	return tagDesc == token ||
		tagLimit == token ||
		tagRegexp == token ||
		tagErrorMessage == token ||
		tagRequire == token ||
		tagDeprecated == token
}

func Parse(tagVal string) map[string]string {
	mp := map[string]string{}
	pos := 0
	state := 0

	token := ""
	val := ""
	tmp := ""
	prior := ""

	for pos < len(tagVal) {
		ch := tagVal[pos : pos+1]
		pos++
		end := pos == len(tagVal)
		switch state {
		case 0: // 找 token
			if end {
				mp[token+ch] = ""
			}
			if ch == "," { // 结束
				if isValidateToken(token) {
					// 合法的令牌，添加到令牌表，继续寻找令牌，状态不变
					mp[token] = ""
					token = ""
					continue
				}
			}
			if ch == ":" {
				if isValidateToken(token) {
					// 合法的令牌
					// 该查找令牌对应的值了
					state = 1
					val = ""
				} else {
					// 非法令牌，继续查找令牌
					token = ""
					val = ""
				}

				continue
			}
			// 完善令牌数据
			token += ch
		case 1: // 找令牌的 Value
			if end {
				if ch != "," {
					val += ch
				}
				mp[token] = val
				continue
			}
			if ch == "," { // 遇到逗号，查找下一个是否令牌
				state = 2
				prior = ch
				tmp = ""
				continue
			}
			val += ch
		case 2: // 查看后面的数据是否令牌
			if end {
				if isValidateToken(tmp) {
					// 是令牌
					mp[token] = val
					mp[tmp] = ""
				} else if isValidateToken(tmp + ch) {
					// 是令牌
					mp[token] = val
					mp[tmp+ch] = ""
				} else {
					// 非令牌
					mp[token] = val + prior + tmp + ch
					prior = ""
				}
				continue
			}
			if ch == "," { // 结束
				if isValidateToken(tmp) {
					// 合法的令牌
					mp[token] = val
					mp[tmp] = ""
					state = 0 // 查找下一个令牌
					prior = ""
					token = ""
				} else {
					val += prior + tmp
					prior = ch
					tmp = ""
				}
				continue
			}
			if ch == ":" { // 该查找值了
				if isValidateToken(tmp) {
					// 合法的令牌
					mp[token] = val
					token = tmp
					state = 1 // 查找令牌的值
					val = ""
					prior = ""
				} else {
					// 非合法令牌
					val += prior + tmp
					prior = ch
					tmp = ""
				}
				continue
			}
			tmp += ch
		}
	}
	return mp
}

// ParseTag 解析 Tag，如果  json tag 设置为 -,则返回 nil,即忽略此字段
func ParseTag(field *reflect.StructField) *ChainTag {
	gTag := &ChainTag{}
	tag := field.Tag

	gTag.FieldName = parseFieldName(&tag, field)
	if gTag.FieldName == "" {
		return nil
	}
	tagVal := tag.Get(TAG)
	if tagVal != "" {
		mp := Parse(tagVal)
		//mp := map[string]string{}
		//ary := strings.Split(tagVal, ",")
		//for _, item := range ary {
		//	p := strings.Index(item, ":")
		//	if p < 0 {
		//		mp[item] = item
		//	} else {
		//		mp[item[:p]] = item[p+1:]
		//	}
		//	//sub := strings.Split(item, "=")
		//	//if len(sub) == 1 {
		//	//	mp[sub[0]] = sub[0]
		//	//}
		//	//if len(sub) == 2 {
		//	//	mp[sub[0]] = sub[1]
		//	//}
		//}

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
