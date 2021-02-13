package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
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
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template //解析后的模板对象
	funcMap       template.FuncMap   //自定义函数映射表
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
func (group *RouterGroup) Group(prefix string) *RouterGroup {
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
	for _, group := range engine.groups {
		if strings.HasPrefix(c.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}
	c.engine = engine //初始化上下文的Engine指针
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//创建静态文件处理器
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//绝对路径,由路由分组前缀和相对路径组成
	absolutePath := path.Join(group.prefix, relativePath)
	//文件服务器:会将请求的URL路径中"绝对路径"去除后,交由fs文件处理器处理
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	//返回判断文件是否存在并执行文件系统服务的函数
	return func(c *Context) {
		file := c.GetParam("filepath") //获取文件路径
		//判断文件是否能打开(存在)
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req) //执行服务器接口函数
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

//添加自定义模板渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载HTML模板
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must( //模板初始化
		template.New(""). //新建匿名模板
					Funcs(engine.funcMap). //加载添加的自定义函数
					ParseGlob(pattern))    //解析模板文件
}
