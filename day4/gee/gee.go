package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)

//路由分组结构体
type RouterGroup struct {
	prefix      string        //分组的前缀
	middlewares []HandlerFunc //中间件
	engine      *Engine       //所属Engine
}

//框架主体结构体
type Engine struct {
	*RouterGroup                //默认的路由分组,未分组的路由都加入该分组
	router       *router        //路由
	groups       []*RouterGroup //存储所有路由分组
}

//框架主体构造函数
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	//构建根路由分组,其engine指向框架主体
	engine.RouterGroup = &RouterGroup{engine: engine}
	//将路由分组记录到路由分组集
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//由根路由分组创建新的路由分组
func (group *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := group.engine //共享所指的框架主体
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	//添加到分组集
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp //完整的路由为分组前缀和当前添加的路径部分
	log.Printf("Route %4s - %s\n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler) //添加路由
}

//GET路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//POST路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
