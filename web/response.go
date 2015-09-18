/**
 * Copyright 2015 @press soft.
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

// 输出JSON
func (this *response) JsonOutput(v interface{}) {
	this.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	b, err := json.Marshal(v)
	if err != nil {
		str := fmt.Sprintf(`{"error":"%s"}`,
			strings.Replace(err.Error(), "\"", "\\\"", -1))
		this.ResponseWriter.Write([]byte(str))
	} else {
		this.ResponseWriter.Write(b)
		//this.Write([]byte(b))
	}
}
