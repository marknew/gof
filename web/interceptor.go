package web

import (
	"errors"
	"fmt"
	"github.com/jrsix/gof"
	"github.com/jrsix/gof/log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	HandleDefaultHttpExcept func(*Context, error)
	HandleHttpBeforePrint   func(*Context) bool
	HandleHttpAfterPrint    func(*Context)
	_                       http.Handler = new(Interceptor)
)

//Http请求处理代理
type Interceptor struct {
	_app gof.App
	//执行请求
	_execute RequestHandler
	//请求之前发生;返回false,则终止运行
	Before func(*Context) bool
	After  func(*Context)
	Except func(*Context, error)
}

func NewInterceptor(app gof.App, f RequestHandler) *Interceptor {
	return &Interceptor{
		_app:     app,
		_execute: f,
	}
}

func (this *Interceptor) handle(app gof.App, w http.ResponseWriter, r *http.Request, handler RequestHandler) {
	// proxy response writer
	//w := NewRespProxyWriter(w)
	ctx := NewContext(app, w, r)

	//todo: panic可以抛出任意对象，所以recover()返回一个interface{}
	if this.Except != nil {
		defer func() {
			if err := recover(); err != nil {
				this.Except(ctx, errors.New(fmt.Sprintf("%s", err)))
			}
		}()
	}

	if this.Before != nil {
		if !this.Before(ctx) {
			return
		}
	}
	if handler != nil {
		handler(ctx)
	}

	if this.After != nil {
		this.After(ctx)
	}
}

func (this *Interceptor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if this._app == nil {
		log.Fatalln("Please use web.NewInterceptor(gof.App) to initialize!")
		os.Exit(1)
	}
	this.handle(this._app, w, r, this._execute)
}

func (this *Interceptor) For(handle RequestHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		this.handle(this._app, w, r, handle)
	})
}

func init() {
	HandleDefaultHttpExcept = func(ctx *Context, err error) {
		_, f, line, _ := runtime.Caller(1)
		var w = ctx.Response
		var header http.Header = w.Header()
		header.Add("Content-Type", "text/html")
		w.WriteHeader(500)

		var part1 string = `<html><head><title>Exception - GOF</title>
				<meta charset="utf-8"/>
				<style>
				body{background:#FFF;font-size:100%;color:#333;margin:0 2%;}
        h1{color:red;font-size:28px;border-bottom:solid 1px #ddd;line-height:80px;}
        div.except-panel p{margin:20px 0;}
        div.except-panel div.summary{}
        div.except-panel p.message{font-size:24px;}
        div.except-panel p.contact{color:#666;font-size:18px;}
        div.except-panel p.stack{padding-top:30px;}
        div.except-panel p.stack em{font-size:18px;font-style: normal;}
        div.except-panel pre{font-family: Sans,Arail;
            border:solid 1px #ddd;padding:20px;
            font-size:16px;background:#F5F5F5;
            line-height: 150%;color:#888;}
        div.except-panel .hidden{display:none;}
			</style>
        </head>
        <body>`

		var html string = fmt.Sprintf(`
				<h1>系统异常：%s</h1>
				<div class="except-panel">
					<div class="summary">
						<p class="message">Source：%s&nbsp;&nbsp;Line:%d</p>
						<p class="contact">请联系管理员或 <a href="/">回到首页</a></p>
					</div>
					<p class="stack">
						<em>堆栈信息：</em><br/>
						<pre>
							%s
						</pre>
					</p>
				</div>
		</body>
		</html>
		`, err.Error(), f, line, debug.Stack())

		w.Write([]byte(part1 + html))
	}

	HandleHttpBeforePrint = func(ctx *Context) bool {
		r := ctx.Request
		fmt.Println("[Request] ", time.Now().Format("2006-01-02 15:04:05"), ": URL:", r.RequestURI)
		for k, v := range r.Header {
			fmt.Println(k, ":", v)
		}
		if r.Method == "POST" {
			r.ParseForm()
		}
		for k, v := range r.Form {
			fmt.Println("form", k, ":", v)
		}
		return true
	}

	HandleHttpAfterPrint = func(ctx *Context) {
		w := ctx.Response
		proxy, ok := w.ResponseWriter.(*ResponseProxyWriter)

		if !ok {
			fmt.Println("[Response] convert error")
			return
		}
		fmt.Println("[Respose]\n" + string(proxy.Output))
	}
}
