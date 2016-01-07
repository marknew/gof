/**
 * Copyright 2015 @ presssfot.com
 * name : template.go
 * author : mark zhang
 * date : -- :
 * description :
 * history :
 */

package gof

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Template
type Template struct {
	Init func(*TemplateDataMap)
}

// the data map for template
type TemplateDataMap map[string]interface{}

//type FuncMap template.FuncMap

func (this TemplateDataMap) Add(key string, v interface{}) {
	this[key] = v
}

func (this TemplateDataMap) Del(key string) {
	delete(this, key)
}

// execute template
func (this *Template) ExecuteWithFunc(w io.Writer, funcMap template.FuncMap, dataMap TemplateDataMap,
	tplPath ...string) error {

	t := template.New("-")

	if funcMap != nil {
		t = t.Funcs(funcMap)
	}

	t, err := t.ParseFiles(tplPath...)
	if err != nil {
		return this.handleError(w, err)
	}

	if this.Init != nil {
		if dataMap == nil {
			dataMap = TemplateDataMap{}
		}
		this.Init(&dataMap)
	}

	err = t.Execute(w, dataMap)

	return this.handleError(w, err)
}

// func (this *Template) Execute(w io.Writer, dataMap TemplateDataMap, tplPath ...string) error {

// 	t, err := template.ParseFiles(tplPath...)
// 	if err != nil {
// 		return this.handleError(w, err)
// 	}

// 	if this.Init != nil {
// 		if dataMap == nil {
// 			dataMap = TemplateDataMap{}
// 		}
// 		this.Init(&dataMap)
// 	}
// 	//follow is write to file for testing
// 	// for _, path := range tplPath {
// 	// 	if strings.HasSuffix(path, "index.html") {
// 	// 		localFile, _ := os.Create("/Users/mark/tmp.html")
// 	// 		defer localFile.Close()
// 	// 		writer := bufio.NewWriter(localFile)

// 	// 		t.Execute(writer, dataMap)
// 	// 	}
// 	// }

// 	err = t.Execute(w, dataMap)
// 	return this.handleError(w, err)
// }

func (this *Template) Execute(w io.Writer, dataMap TemplateDataMap, tplPath ...string) error {

	t, err := template.ParseFiles(tplPath...)
	if err != nil {
		return this.handleError(w, err)
	}

	if this.Init != nil {
		if dataMap == nil {
			dataMap = TemplateDataMap{}
		}
		this.Init(&dataMap)
	}
	//follow is write to file for testing
	for _, path := range tplPath {
		if strings.HasSuffix(path, "sysuserprofile.html") {
			localFile, _ := os.Create("/Users/mark/tmp.html")
			defer localFile.Close()
			writer := bufio.NewWriter(localFile)

			newbytes := bytes.NewBufferString("")
			if err = t.Execute(newbytes, dataMap); err != nil {
				fmt.Println(err.Error())
				return this.handleError(w, err)
			}

			if tplcontent, errbyte := ioutil.ReadAll(newbytes); errbyte != nil {
				return errbyte
			} else {
				tplcontent = []byte(strings.Replace(string(tplcontent), "%2f", "/", -1))
				writer.Write([]byte(template.HTML(tplcontent)))
				writer.Flush()

			}
		}
	}

	newbytes := bytes.NewBufferString("")
	if err = t.Execute(newbytes, dataMap); err != nil {
		return this.handleError(w, err)
	}
	if tplcontent, errbyte := ioutil.ReadAll(newbytes); errbyte != nil {
		return err
	} else {
		if rsp, ok := w.(http.ResponseWriter); ok {
			rsp.Header().Add("Content-Type", "text/html; charset=utf-8")
			//url /会变成%2f
			tplcontent = []byte(strings.Replace(string(tplcontent), "%2f", "/", -1))
			rsp.Write([]byte(template.HTML(tplcontent)))
		}
	}
	//t.Execute(writer, dataMap)
	return nil
	//替代原始写法
}

func (this *Template) handleError(w io.Writer, err error) error {
	if err != nil {
		if rsp, ok := w.(http.ResponseWriter); ok {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}
	}
	return err
}
