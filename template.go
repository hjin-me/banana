package banana

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/russross/blackfriday"
)

var (
	t5xx            *template.Template = template.Must(template.New("5xx").Parse("Template file error"))
	ErrTplNotExist  error              = errors.New("tpl not exist")
	ErrTplParseFail error              = errors.New("tpl parse fail")
)

type themeConf map[string]string

func loadThemeConf(filename string) (themeConf, error) {
	x := make(themeConf)
	err := Config(filename, &x)
	dir := filepath.Dir(filename)
	for k, v := range x {
		x[k] = filepath.Join(dir, v)
	}
	return x, err
}

func LoadTheme(dir string) (t *template.Template, err error) {
	fmConf := fisLoadMap(dir)

	themeName := "t:banana"
	t = template.New(themeName)
	funcMaps := template.FuncMap{
		"md":      markDowner,
		"require": fisRequire,
		"uri":     fisURI,
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
			return t, err
		}
		t, err = t.New(name).Parse(string(b))
		if err != nil {
			log.Println("parse tpl:", err, name)
			return t, err
		}
	}
	return
}

func LoadTpl(path string) (*template.Template, string, error) {
	var err error

	funcMaps := template.FuncMap{
		"md":      markDowner,
		"require": fisRequire,
	}
	tc := template.New(path).Funcs(funcMaps)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("load tpl failed:", err)
		return t5xx, "", ErrTplNotExist
	}
	s := string(b)

	tc, err = tc.Parse(s)
	if err != nil {
		log.Println("load tpl failed:", err)
		return t5xx, "", ErrTplParseFail
	}
	return tc, path, nil
}

func Render5xx(w io.Writer, err error) {
	t5xx.Execute(w, err.Error())
}

func Render(w io.Writer, t *template.Template, name string, data interface{}) {
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		panic(err)
	}
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}
