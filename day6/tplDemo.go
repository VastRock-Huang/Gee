package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", testTpl)
	log.Fatal(http.ListenAndServe(":9990", nil))
}

func testTpl(w http.ResponseWriter, req *http.Request) {
	tpl := template.New("demo") //创建一个模板对象
	t := template.Must(         //Must用于初始化解析后的模板
		//解析模板文件
		tpl.Funcs(
			template.FuncMap{
				"fun": func(arg string) (string, error) {
					return "hello " + arg, nil
				},
			}).Parse("<html><body>{{fun .}}</body></html>"))
	//执行模板渲染
	t.Execute(w, "test")
}
