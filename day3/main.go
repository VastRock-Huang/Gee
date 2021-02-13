package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gee!</h1>")
	})

	r.GET("/hello", func(ctx *gee.Context) {
		ctx.String(http.StatusOK, "Hello %s, you're at %s\n",
			ctx.Query("name"), ctx.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n",
			c.GetParam("name"), c.Path)
	})

	r.GET("/hello/:file", func(c *gee.Context) {
		c.String(http.StatusOK, "hello !bin, you're at %s\n",
			c.Path)
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.GetParam("filepath")})
	})

	r.Run(":9999")
}

