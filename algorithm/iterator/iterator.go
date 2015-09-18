/**
 * Copyright 2015 @press soft.
 * name : iterator
 * author : mark zhang
 * date : -- :
 * description : 迭代器
 * history :
 */
package iterator

// 处理单个对象
type WalkFunc func(v interface{}, level int)

// 迭代时满足的条件
type Condition func(v, v1 interface{}) bool

// 迭代栏目,start为开始前执行函数,over为结束迭代执行函数
func Walk(categories []interface{}, v interface{}, c Condition,
	start WalkFunc, over WalkFunc, level int) {
	start(v, level)
	for _, v1 := range categories {
		if c(v, v1) {
			Walk(categories, v1, c, start, over, level+1)
		}
	}
	if over != nil {
		over(v, level)
	}
}
