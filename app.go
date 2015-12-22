/**
 * Copyright 2014 @press soft.
 * name : app1.go
 * author : mark zhang
 * date : 2015-04-27 20:43:
 * description :
 * history :
 */
package gof

import (
	"github.com/mark/gof/db"
	"github.com/mark/gof/log"
)

// 应用当前的上下文
var CurrentApp App

type App interface {
	// Provided db access
	Db() db.Connector
	// Return a Wrapper for GoLang template.
	Template() *Template
	// Return application configs.
	Config() *Config
	// Storage
	Storage() Storage
	// Return a logger
	Log() log.ILogger
	// Application is running debug mode
	Debug() bool
}
