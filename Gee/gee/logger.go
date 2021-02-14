package gee
//日志中间件部分
import (
	"log"
	"time"
)

//返回日志中间件函数
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		//输出记录处理时间
		log.Printf("[%d] %s in %v", c.StatusCode,
			c.Request.RequestURI, time.Since(t))
	}
}
