package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int
}

func FormatAsData(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	//添加自定义模板函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsData": FormatAsData,
	})
	//加载模板目录
	r.LoadHTMLGlob("templates/*")
	//设置静态资源
	r.Static("/assets", "./static")
	stu1 := &student{Name: "Geektutu", Age: 15}
	stu2 := &student{Name: "hh", Age: 20}
	r.GET("/", func(c *gee.Context) {
		//使用css.tmpl渲染
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		//使用arr.tmpl渲染
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/data", func(c *gee.Context) {
		//使用custom_func.tmpl渲染
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
		})
	})
	r.Run(":9999")
}
