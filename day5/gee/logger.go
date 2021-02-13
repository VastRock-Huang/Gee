package gee

import (
	"log"
	"time"
)

//记录处理时间
func Logger() HandlerFunc {
	return func(c *Context) {
		t:=time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode,
			c.Req.RequestURI, time.Since(t))
	}
}
