package banana

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"
)

type Context interface {
	context.Context
	Res() http.ResponseWriter
	Req() *http.Request
	Params() map[string]string
	Output(interface{}, string)
	Json(interface{})
	Tpl(string, interface{})
}

type httpContext struct {
	context.Context
	w http.ResponseWriter
	r *http.Request
	p map[string]string
}

func WithHttp(parent context.Context, w http.ResponseWriter, r *http.Request, p map[string]string) Context {
	return &httpContext{parent, w, r, p}
}

func (c *httpContext) Res() http.ResponseWriter {
	return c.w
}

func (c *httpContext) Req() *http.Request {
	return c.r
}

func (c *httpContext) Params() map[string]string {
	return c.p
}

func (c *httpContext) Output(data interface{}, contentType string) {
	res := c.Res()

	select {
	case <-c.Done():
		log.Println("request timeout", c.Err())
	default:
		h := res.Header()
		h.Add("content-type", contentType)
		fmt.Fprintf(res, "%s", data)
	}
}

func (c *httpContext) Json(data interface{}) {
	res := c.Res()
	select {
	case <-c.Done():
		log.Println("request timeout", c.Err())
	default:
		h := res.Header()
		h.Add("content-type", "application/json")

		str, _ := json.Marshal(data)
		fmt.Fprintf(res, "%s", str)
	}
}

func (c *httpContext) Tpl(path string, data interface{}) {
	cfg, ok := c.Value("cfg").(AppCfg)
	if !ok {
		log.Println("configuration is not ok")
		c.Output("configuration err", "text/plain")
		return
	}
	tpl := filepath.Join(cfg.Env.Tpl, path)
	select {
	case <-c.Done():
		log.Println("request timeout", c.Err())
	default:
		h := c.Res().Header()
		h.Add("content-type", "text/html")
		Render(c.Res(), tpl, data)
	}
}
