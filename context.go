package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	StatusCode int

	Method string
	Path   string
	Params map[string]string

	handlers []HandlerFunc
	index    int

	engine *Engine
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Method:  r.Method,
		Path:    r.URL.Path,
		index:   -1,
	}
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(status int) {
	c.StatusCode = status
	c.Writer.WriteHeader(status)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(status int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(status)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context) JSON(status int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(status)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(status int, data []byte) {
	c.Status(status)
	_, _ = c.Writer.Write(data)
}

func (c *Context) HTML(status int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(status)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
