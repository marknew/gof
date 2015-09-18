/**
 * Copyright 2015 @press soft.
 * name : string.go
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */
package util

import (
	"crypto/rand"
)

const letterStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

//随机字符号串
func RandString(n int) string {
	lsLen := len(letterStr)
	var runes = make([]byte, n)
	rand.Read(runes)
	for i, b := range runes {
		runes[i] = letterStr[b%byte(lsLen)]
	}
	return string(runes)
}
