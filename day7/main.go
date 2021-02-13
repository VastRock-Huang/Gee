package main

import (
	"gee"
	"net/http"
)

func main() {
	r:=gee.Default()	//带Logger和Recovery的框架
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK,"Hello Gee\n")
	})
	r.GET("/panic", func(c *gee.Context) {
		names:=[]string{"geektutu"}
		c.String(http.StatusOK, names[100])		//数组越界会引发panic
	})
	r.GET("/json", func(c *gee.Context) {
		c.JSON(http.StatusOK,gee.H{
			"name":"gee",
		})
	})
	r.Run(":9999")
}