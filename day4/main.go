package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.NewGroup("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee in v1</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "Hello %s, you're at %s\n",
				c.Query("name"), c.Path)
		})
	}
	v2 := r.NewGroup("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s in v2\n",
				c.GetParam("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
		v3 := v2.NewGroup("/v3")
		{
			v3.GET("/print", func(c *gee.Context) {
				c.HTML(http.StatusOK, "<p>It's v3 print</p>")
			})
		}
	}
	r.Run(":9999")
}
