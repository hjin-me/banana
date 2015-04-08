package banana

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"sync"

	"github.com/russross/blackfriday"
)

var (
	ErrTplNotExist  = errors.New("tpl not exist")
	ErrTplParseFail = errors.New("tpl parse fail")

	themeCache = make(map[string]*tCache)
	once       sync.Once
)

type tCache struct {
	once sync.Once
	t    *template.Template
}

func LoadTheme(dir string) (t *template.Template, err error) {
	cache, ok := themeCache[dir]
	if !ok {
		cache = &tCache{}
		cache.once.Do(func() {
			fmConf := fisLoadMap(dir)

			themeName := "t:banana"
			t = template.New(themeName)
			funcMaps := template.FuncMap{
				"md":   markDowner,
				"link": fisLink,
				"uri":  fisURI,
				"block": func(name string, data interface{}) (ret template.HTML, err error) {
					buf := bytes.NewBuffer([]byte{})
					err = t.ExecuteTemplate(buf, name, data)
					ret = template.HTML(buf.String())
					return
				},
			}
			t.Funcs(funcMaps)

			for name, fr := range fmConf.Res {
				if !fr.IsPage() {
					continue
				}
				b, err := ioutil.ReadFile(fr.URI)
				if err != nil {
					log.Println("load tpl failed:", err)
					return
				}
				t, err = t.New(name).Parse(string(b))
				if err != nil {
					log.Println("parse tpl:", err, name)
					return
				}
			}
			cache.t = t
			themeCache[dir] = cache
		})
	}

	return cache.t, err
}

func Render(t *template.Template, name string, data interface{}) (bf bytes.Buffer, err error) {
	err = t.ExecuteTemplate(&bf, name, data)
	return
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
