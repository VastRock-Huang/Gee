package main

import (
	"fmt"
	"log"
	"net/http"
)

//Engine是一个对所有请求的统一句柄
type Engine struct {}

//注意函数名不能修改, 要对应
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w,"URL.Path= %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header{
			fmt.Fprintf(w,"Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main()  {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
	//Engine类型会强制转换为Handler接口类型,调用其中的ServeHTTP方法
}

