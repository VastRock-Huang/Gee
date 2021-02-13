package gee

import (
	"log"
	"net/http"
)

//路由匹配部分

//路由结构体
type router struct {
	handlers map[string]HandlerFunc		//路由映射表
}

//路由构造函数
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

//添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc)  {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

//路由处理
func (r *router) handle(c *Context)  {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}