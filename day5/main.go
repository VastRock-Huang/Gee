package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		//此处使用c.Fail是用于验证中间件执行是否成功,会直接返回报文
		c.Fail(http.StatusInternalServerError, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode,
			c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	v1 := r.NewGroup("/v1")
	v1.Use(gee.Logger())
	v1.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>hello gee</h1>")
	})
	v2 := r.NewGroup("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n",
				c.GetParam("name"), c.Path)
		})
	}
	if err := r.Run(":9999"); err != nil {
		fmt.Println(err)
	}
}
