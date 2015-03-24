package banana

import (
	"fmt"
	"html/template"
	"path/filepath"
)

var (
	extMap = make(map[string]string)
)

func fisRequire(args ...interface{}) template.HTML {

	extMap[".js"] = "js"
	extMap[".css"] = "css"
	extMap[".less"] = "css"
	path, ok := args[0].(string)
	if !ok {
		panic("args[0] is not string")
	}

	tp := extMap[filepath.Ext(path)]
	if !ok {
		panic("unknown type [" + path + "]")
	}
	s := ""
	switch tp {
	case "js":
		s = fmt.Sprintf("<script src=\"%s\"></script>", path)
	case "css":
		s = fmt.Sprintf("<link type=\"text/css\" rel=\"stylesheet\" href=\"%s\">", path)
	}

	return template.HTML(s)
}
