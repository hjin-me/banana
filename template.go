package banana

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
)

var (
	t5xx *template.Template = template.Must(template.New("5xx").Parse("Template file error"))
)

func Load5xx() (*template.Template, string) {
	log.Println("load new 5xx tpl")

	return t5xx, "5xx"
}

func LoadTpl(path string) (*template.Template, string) {
	var err error
	tc := template.New(path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("load tpl failed:", err)
		return Load5xx()
	}
	s := string(b)

	tc, err = tc.Parse(s)
	if err != nil {
		log.Println("load tpl failed:", err)
		return Load5xx()
	}
	return tc, path
}

func Render(w io.Writer, path string, data interface{}) {
	t, name := LoadTpl(path)
	t.ExecuteTemplate(w, name, data)
}
