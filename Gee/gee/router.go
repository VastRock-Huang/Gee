package gee
//路由匹配部分
import (
	"net/http"
	"strings"
)

//路由匹配结构体
type router struct {
	roots    map[string]*node	//前缀树根结点,以请求方法作键名
	handlers map[string]HandlerFunc	//路由处理函数映射表
}

//router构造函数
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

//解析模式字符串或具体路由字符串成前缀字符串集(将路径按"/"分割)
func parsePattern(pattern string) []string {
	//按"/"分割的字符串数组
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {	//将不为空字符串的添加到前缀
			parts = append(parts, item)
			//遇到通配符则退出, 不考虑后续前缀
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)	//解析路径获得前缀
	key := method + "-" + pattern	//路由映射表键名
	//若对应请求的方法不存在则构建根结点
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	//添加路由到前缀树
	r.roots[method].insert(pattern, parts, 0)
	//添加路由处理函数
	r.handlers[key] = handler
}

//根据具体路由返回对应前缀树结点和对应参数映射表
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)	//解析路径获得前缀
	params := make(map[string]string)
	root, ok := r.roots[method]		//获得请求方法的前缀树根结点
	//结点不存在则返回nil
	if !ok {
		return nil, nil
	}
	//寻找一个和前缀匹配的结点
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)	//解析匹配结点的模式字符串
		for idx, part := range parts {
			//若为动态路由
			if part[0] == ':' {
				//将动态路由名和路径中对应的参数值作为键值对添加到params映射中
				params[part[1:]] = searchParts[idx]
			} else if part[0] == '*' {
				//将通配参数名和路径中对应参数值添加到param映射
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
				break
			}
		}
	}
	return n, params	//返回匹配结点和映射表
}

//获取请求方法对应的前缀树的所有结点
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)	//递归遍历前缀树的结点并添加
	return nodes
}

//路由处理
func (r *router) handle(c *Context) {
	//获取路由的前缀树结点和对应参数表
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		//根据结点的模式字符串确定待匹配路由
		key := c.Method + "-" + n.pattern
		//将路由处理函数添加到上下文的处理函数集中
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		//将404函数添加到上下文的处理函数集中
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND :%s\n",
				c.Path)
		})
	}
	c.Next()	//开始执行处理函数
}
