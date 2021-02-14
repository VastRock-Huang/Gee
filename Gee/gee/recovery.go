package gee

//错误恢复中间件
import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//输出调用栈信息
func trace(message string) string {
	var pcs [32]uintptr             //指针
	n := runtime.Callers(3, pcs[:]) //获取调用函数指针(跳过前3个)
	var str strings.Builder         //快速字符串拼接
	str.WriteString(message)        //拼接字符串
	str.WriteString("\nTraceback:")
	for _, pc := range pcs[:n] { //遍历函数指针
		fn := runtime.FuncForPC(pc)   //返回函数指针对应函数
		file, line := fn.FileLine(pc) //文件及行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String() //转字符串
}

//返回错误恢复中间件函数
func Recovery() HandlerFunc {
	//返回中间件函数
	return func(c *Context) {
		//错误处理
		defer func() {
			//恢复错误并输出错误信息
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				//终止请求的处理
				c.Fail(http.StatusInternalServerError,
					"Internal Server Error")
			}
		}()
		c.Next() //用于跳转至下一中间件(可省略)
	}
}
