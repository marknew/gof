/**
 * Copyright 2015 @ presssfot.com
 * name : session.go
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */
package web

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mark/gof"
	"github.com/mark/gof/util"
)

const (
	defaultSessionMaxAge int64  = 3600 * 12
	sessionCookieName    string = "gof_SessionId"
)

var (
	sessionStorage gof.Storage
)

type Session struct {
	_sessionId string
	_rsp       http.ResponseWriter
	_data      map[string]interface{}
	_storage   gof.Storage
	_maxAge    int64
}

func getSession(w http.ResponseWriter, storage gof.Storage, sessionId string) *Session {
	s := &Session{
		_sessionId: sessionId,
		_rsp:       w,
		_data:      make(map[string]interface{}),
		_storage:   storage,
		_maxAge:    defaultSessionMaxAge,
	}
	s._storage.Get(getSessionId(s._sessionId), &s._data)
	return s
}

func newSession(w http.ResponseWriter, storage gof.Storage) *Session {
	id := newSessionId()
	return &Session{
		_sessionId: id,
		_rsp:       w,
		_storage:   storage,
		_maxAge:    defaultSessionMaxAge,
	}
}

func (this *Session) chkInit() {
	if this._data == nil {
		this._data = make(map[string]interface{})
	}
}

// 获取会话编号
func (this *Session) GetSessionId() string {
	return this._sessionId
}

// 获取值
func (this *Session) Get(key string) interface{} {
	if this._data != nil {
		if v, ok := this._data[key]; ok {
			return v
		}
	}
	return nil
}

// 设置键值
func (this *Session) Set(key string, v interface{}) {
	this.chkInit()
	this._data[key] = v
}

// 移除键
func (this *Session) Remove(key string) bool {
	if _, exists := this._data[key]; exists {
		delete(this._data, key)
		return true
	}
	return false
}

// 使用指定的会话代替当前会话
func (this *Session) UseInstead(sessionId string) {
	this._sessionId = sessionId
	this._storage.Get(getSessionId(this._sessionId), &this._data)
	this.flushToClient()
}

// 销毁会话
func (this *Session) Destroy() {
	this._storage.Del(getSessionId(this._sessionId))
	this.SetMaxAge(-this._maxAge)
	this.flushToClient()
}

// 保存会话
func (this *Session) Save() error {
	if this._data == nil {
		return nil
	}

	err := this._storage.SetExpire(getSessionId(this._sessionId), &this._data, this._maxAge)
	if err == nil {
		this.flushToClient()
	}
	return err
}

// 设置过期秒数
func (this *Session) SetMaxAge(seconds int64) {
	this._maxAge = seconds
}

//存储到客户端
func (this *Session) flushToClient() {
	d := time.Duration(this._maxAge * 1e9)
	expires := time.Now().Local().Add(d)
	ck := &http.Cookie{
		Name:     sessionCookieName,
		Value:    this._sessionId,
		Path:     "/",
		HttpOnly: true,
		Expires:  expires,
	}
	http.SetCookie(this._rsp, ck)
}

func init() {
	gob.Register(make(map[string]interface{})) // register session type for gob.
}

func getSessionId(id string) string {
	return "gof:web:session:" + id
}

func newSessionId() string {
	dt := time.Now()
	randStr := util.RandString(4)
	return fmt.Sprintf("%s%d%d", randStr, dt.Second(), dt.Nanosecond())
}

// Set global session storage
func SetSessionStorage(s gof.Storage) {
	sessionStorage = s
}

func parseSession(rsp http.ResponseWriter, r *http.Request, cookieName string, sto gof.Storage) *Session {
	var s *Session
	ck, err := r.Cookie(cookieName)
	if sto == nil {
		log.Fatalln("session storage is nil")
	}
	if err == nil {
		s = getSession(rsp, sto, ck.Value)
	} else {
		s = newSession(rsp, sto)
	}
	return s
}

// Session adapter for http context
func SessionAdapter(rsp http.ResponseWriter, r *http.Request) *Session {
	return parseSession(rsp, r, sessionCookieName, sessionStorage)
}
