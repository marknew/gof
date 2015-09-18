/**
 * Copyright 2015 @press soft.
 * name : fmt
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */

package fmt

func BoolString(b bool, trueVal, falseVal string) string {
	if b {
		return trueVal
	}
	return falseVal
}

func BoolInt(b bool, v, v1 int) int {
	if b {
		return v
	}
	return v1
}
