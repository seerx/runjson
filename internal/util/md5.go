package util

import (
	"crypto/md5"
	"fmt"
)

//MD5 MD5 加密
func MD5(vals ...string) string {
	val := ""
	for _, v := range vals {
		val += v
	}
	data := []byte(val)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
