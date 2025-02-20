package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 将 ResponseWriter 和 Request 封装到 Context 中
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 通过 Context 结构体获取PostForm参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 通过 Context 结构体获取Query参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 通过 Context 结构体设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 通过 Context 结构体设置Header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 通过 Context 结构体返回String类型的响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 通过 Context 结构体返回JSON类型的响应
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 通过 Context 结构体返回Data类型的响应
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 通过 Context 结构体返回HTML类型的响应
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
