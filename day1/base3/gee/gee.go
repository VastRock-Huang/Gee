package gee

import (
	"fmt"
	"net/http"
)

//请求处理函数
type HandlerFunc func(http.ResponseWriter, *http.Request)

//Engine 实现ServeHTTP接口
type Engine struct {
	router map[string]HandlerFunc	//路由表: 方法-路径 : 处理函数
}

//Engine构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

//添加路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc)  {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//添加GET请求的路由
func (engine *Engine) GET(pattern string, handler HandlerFunc)  {
	engine.addRoute("GET", pattern, handler)
}

//添加POST请求的路由
func (engine *Engine) POST(pattern string, handler HandlerFunc)  {
	engine.addRoute("POST", pattern, handler)
}

//运行服务端
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//路由匹配函数
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	key := req.Method + "-" +req.URL.Path	//获取路由的键名
	//路由表中存在则使用相应的处理函数
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {	//不存在则404
		w.WriteHeader(http.StatusNotFound)	//构造响应, 返回错误码
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}