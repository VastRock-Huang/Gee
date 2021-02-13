package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

func FormatAsData(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.Default() //带Logger和Recovery的框架
	r.Static("/assets", "./static")
	r.SetFuncMap(template.FuncMap{
		"FormatAsData": FormatAsData,
	})
	r.LoadHTMLGlob("./templates/*")
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/panic", func(c *gee.Context) {
			names := []string{"geektutu"}
			c.String(http.StatusOK, names[100]) //数组越界会引发panic
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/data/:name", func(c *gee.Context) {
			t := time.Now()
			c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
				"title": c.GetParam("name"),
				"now": time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
					t.Second(), t.Nanosecond(), time.Local),
			})
		})
	}
	r.Run(":9999")
}
