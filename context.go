package banana

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Output(interface{}, string) error
	Json(interface{}) error
	Tpl(string, interface{}) error
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

func (c *httpContext) Output(data interface{}, contentType string) (err error) {
	res := c.Res()

	select {
	case <-c.Done():
		err = c.Err()
	default:
		h := res.Header()
		h.Add("content-type", contentType)
		fmt.Fprintf(res, "%s", data)
	}
	return
}

func (c *httpContext) Json(data interface{}) (err error) {
	res := c.Res()
	select {
	case <-c.Done():
		err = c.Err()
	default:
		h := res.Header()
		h.Add("content-type", "application/json")

		str, _ := json.Marshal(data)
		fmt.Fprintf(res, "%s", str)
	}
	return
}

func (c *httpContext) Tpl(path string, data interface{}) (err error) {
	var bf bytes.Buffer
	ch := make(chan error)
	go func() {
		cfg, ok := c.Value("cfg").(AppCfg)
		if !ok {
			log.Println("configuration is not ok")
			err = errors.New("configuration err")
			return
		}
		themeDir := filepath.Join(cfg.Env.Tpl)
		h := c.Res().Header()
		h.Add("content-type", "text/html")
		t, err := LoadTheme(themeDir)

		select {
		case <-c.Done():
			err = c.Err()
		default:
			if err == nil {
				bf, err = Render(t, path, data)
			}

			ch <- err
		}
	}()
	select {
	case <-c.Done():
		err = c.Err()
	case err = <-ch:
		if err == nil {
			bf.WriteTo(c.Res())
		}
	}
	return
}
