package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

const passSalt = "xulei$xulei"

// Md5Password 密码加密
func Md5Password(pass string) string {
	w := md5.New()

	io.WriteString(w, pass+passSalt)     //将str写入到w中
	return fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
}
