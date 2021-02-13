package gee

import "net/http"

type HandlerFunc func(c * Context)

//Gee主体结构体
type Engine struct {
	router *router //路由
}

//Engine构造函数
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

//添加路由
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc)  {
	engine.router.addRoute(method, pattern, handler)
}

//添加GET路由
func (engine *Engine) GET(pattern string, handler HandlerFunc)  {
	engine.addRoute("GET", pattern, handler)
}

//添加POST路由
func (engine *Engine) POST(pattern string, handler HandlerFunc)  {
	engine.addRoute("POST", pattern, handler)
}

//服务端http.Handler接口函数
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request)  {
	c := newContext(w, req)
	engine.router.handle(c)
}

//运行框架
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
