/**
 * Copyright 2015 @ presssfot.com
 * name : controller.go
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */
package mvc

type Controller interface{}

// Generate controller instance
type ControllerGenerate func() Controller
