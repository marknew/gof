/**
 * Copyright 2015 @ presssfot.com
 * name : response
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var _ http.ResponseWriter = new(response)

type response struct {
	http.ResponseWriter
}

//输出JSON
//处理从JQuery传来的callback
func (this *response) JsonOutput(ctx *Context, v interface{}) {

	this.ResponseWriter.Header().Set("Content-Type", "application/json")
	jsondata, err := json.Marshal(v)
	// fmt.Println(string(jsondata))
	if callback := ctx.Request.FormValue("callback"); callback != "" {
		//this.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		jsondata = []byte(fmt.Sprintf("%s(%s);", callback, jsondata))
	}

	if err != nil {
		str := fmt.Sprintf(`{"error":"%s"}`,
			strings.Replace(err.Error(), "\"", "\\\"", -1))
		this.ResponseWriter.Write([]byte(str))
	} else {
		this.ResponseWriter.Write(jsondata)
	}
}
