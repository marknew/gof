package util

import (
	"fmt"

	"runtime"

	log "github.com/golang/glog"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func HandleError(errs string) {
	var err error
	if r := recover(); r != nil {
		switch r := r.(type) {
		case error:
			err = r
		default:
			err = fmt.Errorf("%v", r)
		}
		log.V(2).Infoln(errs, err.Error())
	}
}
func IIF(b bool, d1, d2 interface{}) interface{} {
	if b {
		return d1
	} else {
		return d2
	}
}

/*
\033[0m 关闭所有属性
\033[1m 设置高亮度
\033[4m 下划线
\033[5m 闪烁
\033[7m 反显
\033[8m 消隐
\033[30m 至 \33[37m 设置前景色
\033[40m 至 \33[47m 设置背景色
\033[nA 光标上移n行
\033[nB 光标下移n行
\033[nC 光标右移n行
\033[nD 光标左移n行
\033[y;xH设置光标位置
\033[2J 清屏
\033[K 清除从光标到行尾的内容
\033[s 保存光标位置
\033[u 恢复光标位置
\033[?25l 隐藏光标
\033[?25h 显示光标<br>

40----49

40:黑

41:深红

42:绿

43:黄色

44:蓝色

45:紫色

46:深绿

47:白色

字颜色:30-----------39

30:黑

31:红

32:绿

33:黄

34:蓝色

35:紫色

36:深绿

37:白色
\033[32;1m我被变成了蓝色，\033[0m我是原来的颜色
*/
func Colorize(text string, status string) string {
	if runtime.GOOS == "windows" {
		return text
	}
	out := ""
	switch status {
	case "succ":
		out = "\033[32;1m" // Blue
	case "fail":
		out = "\033[31;1m" // Red
	case "warn":
		out = "\033[33;1m" // Yellow
	case "note":
		out = "\033[34;1m" // Green
	default:
		out = "\033[0m" // Default
	}
	return out + text + "\033[0m"
}
