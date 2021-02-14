package gee

//上下文部分
import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

//上下文结构体
type Context struct {
	Writer     http.ResponseWriter //响应
	Request    *http.Request       //请求
	Path       string              //请求的URL路径
	Method     string              //请求的方法
	StatusCode int                 //响应的状态码
	Params     map[string]string   //动态路由参数表
	handlers   []HandlerFunc       //处理函数集
	index      int                 //处理函数索引
	engine     *Engine             //框架主体指针
}

//Context构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: req,
		Path:    req.URL.Path,
		Method:  req.Method,
		index:   -1,
	}
}

//依次执行中间件
func (c *Context) Next() {
	c.index++
	for s := len(c.handlers); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//执行失败中断中间件执行
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)       //将处理函数索引设置为终止值
	c.JSON(code, H{"message": err}) //返回错误信息
}

//获取动态路由的参数
func (c *Context) GetParam(part string) string {
	return c.Params[part]
}

//根据键名key获取请求的表单中对应的键值的第一个
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

//根据键名key获取URL查询字符串中对应的键值的第一个
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

//添加响应的首部
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//设置响应状态码
func (c *Context) SetStatus(code int) {
	c.Writer.WriteHeader(code)
}

//构造文本类型响应
func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

//构造JSON类型响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err) //JSON编码失败,则引发错误
	}
}

//构造普通数据响应
func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	c.Writer.Write(data)
}

//构造超文本类型响应并进行渲染
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	//执行模板渲染
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		panic(err) //HTML模板渲染失败,则引发错误
	}
}
