package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr //指针
	//获取调用栈指针(跳过前3个)
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder  //快速字符串拼接
	str.WriteString(message) //拼接字符串
	str.WriteString("\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)		//返回栈指针对应函数
		file, line := fn.FileLine(pc)	//文件及行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()	//转字符串
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError,
					"Internal Server Error")
			}
		}()
		c.Next()
	}
}
