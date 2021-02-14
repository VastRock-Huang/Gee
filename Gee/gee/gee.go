package gee

//框架主体部分
import (
	"html/template"
	"net/http"
	"strings"
)

//处理函数
type HandlerFunc func(c *Context)

//框架主体结构体
type Engine struct {
	*RouterGroup                     //默认的路由分组,未分组的路由都加入该分组
	router        *router            //路由
	groups        []*RouterGroup     //存储所有路由分组
	htmlTemplates *template.Template //解析后的模板对象指针
	funcMap       template.FuncMap   //自定义模板渲染函数映射表
}

//Engine构造函数
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	//构建根路由分组,其engine指向框架主体
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	//将路由分组记录到路由分组集
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//默认框架
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery()) //带日志记录和异常恢复
	return engine
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
			ParseGlob(pattern)) //解析模板文件
}

//服务端http.Handler接口函数
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)               //创建针对该请求的上下文
	for _, group := range engine.groups { //遍历所有路由分组筛选中间件
		//若请求的URL的前缀与路由分组相同
		if c.Path == group.prefix || strings.HasPrefix(c.Path, group.prefix+"/") {
			//添加中间件函数
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}
	c.engine = engine       //初始化上下文的Engine指针
	engine.router.handle(c) //执行路由处理
}

//运行框架
func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
