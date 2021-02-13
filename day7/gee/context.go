package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string
	handlers   []HandlerFunc //处理函数集(存放中间件和路由处理函数)
	index      int           //处理函数索引
	engine     *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

//依次执行中间件
func (c *Context) Next() {
	c.index++
	for l := len(c.handlers); c.index < l; c.index++ {
		c.handlers[c.index](c)
	}
}

//执行失败
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers) //将处理函数索引直接设置为终止值, 停止后续处理函数执行
	c.JSON(code, H{"message": err})
}

func (c *Context) GetParam(part string) string {
	return c.Params[part]
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetStatus(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err == nil {
		panic(err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatus(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		//c.Fail(500, err.Error())
		panic(err)
	}
}
