package gee

//路由分组部分
import (
	"log"
	"net/http"
	"path"
)

//路由分组结构体
type RouterGroup struct {
	prefix      string        //分组的前缀
	middlewares []HandlerFunc //中间件函数集
	engine      *Engine       //框架主体指针
}

//由父路由分组创建新的路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine //共享所指的框架主体
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix, //以父分组前缀构建新前缀
		engine: engine,
	}
	//添加到框架主体的路由分组集中
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	//完整的路由为分组前缀和当前添加的路径部分
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	//添加路由
	group.engine.router.addRoute(method, pattern, handler)
}

//添加GET路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//添加POST路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//向路由分组中添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//创建静态文件处理函数
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//绝对路径,由路由分组前缀和相对路径组成
	absolutePath := path.Join(group.prefix, relativePath)
	//文件服务器:会将请求的URL路径中"绝对路径"去除后,交由fs文件处理器处理
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	//返回判断文件是否存在并执行文件系统服务的函数
	return func(c *Context) {
		//获取文件路径
		file := c.GetParam("filepath")
		//判断文件是否能打开(存在)
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		//执行服务器接口函数
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

//添加静态文件路由
//relativePath是文件的相对路径,root是映射到的项目目录
func (group *RouterGroup) Static(relativePath string, root string) {
	//静态文件处理函数
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	//静态文件的模式字符串
	urlPattern := path.Join(relativePath, "/*filepath")
	//添加路由
	group.GET(urlPattern, handler)
}
