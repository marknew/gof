/**
 * Copyright 2015 @ presssfot.com
 * name : filter.go
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */
package mvc

import (
	"github.com/mark/gof/web"
)

// controller filter
type Filter interface {
	//call it before execute your some business.
	Requesting(*web.Context) bool
	//call it after execute your some business.
	RequestEnd(*web.Context)
}
