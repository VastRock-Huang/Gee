package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//上下文部分

type H map[string]interface{}

//上下文结构体
type Context struct {
	Writer http.ResponseWriter	//回复
	Req *http.Request	//请求
	Path string		//请求的URL路径
	Method string	//请求的方法
	StatusCode int	//状态码
}

//Context构造函数
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
	}
}

//根据键名key获取请求的表单中对应的键值的第一个
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//根据键名key获取URL查询字符串中对应的键值的第一个
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//设置响应状态码
func (c *Context) Status(code int)  {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//添加响应的首部
func (c *Context) SetHeader(key string, value string)  {
	c.Writer.Header().Set(key, value)
}

//构造文本类型响应
func (c *Context) String(code int, format string, values ...interface{})  {
	c.SetHeader("Content-Type","text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//构造JSON类型响应
func (c *Context) JSON(code int, obj interface{})  {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err  != nil {
		//http.Error(c.Writer, err.Error(), 500)
		//http.Error中使用了Header().Set()和WriteHeader(),
		//但encoder.Encode(obj)内部实现相当于调用了Write(),
		//因此http.Error函数在出错时并不能正常工作, 此处参考Gin源码采用了异常处理
		log.Fatal("Error")
		panic(err)
	}
}

//构造普通数据响应
func (c *Context) Data(code int, data []byte)  {
	c.Status(code)
	c.Writer.Write(data)
}

//构造超文本类型响应
func (c *Context) HTML(code int, html string)  {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}