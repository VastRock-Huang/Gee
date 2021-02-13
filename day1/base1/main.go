package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main()  {
	//注册路由和处理函数到默认路由
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	//监听端口并处理, nil表示使用默认路由
	err:= http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//处理函数
func indexHandler(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprintf(w, "URL.Path= %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request){
	io.WriteString(w, "hello world!\n")
	for k, v :=range req.Header{
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
