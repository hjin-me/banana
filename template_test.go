package banana

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadTheme(t *testing.T) {
	themeDir, _ := filepath.Abs("test/views/cp/")
	tpl, err := LoadTheme(themeDir)
	if err != nil {
		t.Error(err)
	}
	var wb bytes.Buffer
	tpl.ExecuteTemplate(&wb, "layout", "nothing")
	html := wb.String()
	t.Log(html)
	t.Skip("temp")
	if !strings.Contains(html, "layout") {
		t.Error("no layout")
	}
	if !strings.Contains(html, "Head") {
		t.Error("no Head")
	}
}

func TestMdFunc(t *testing.T) {
	themeDir, _ := filepath.Abs("test/output/")
	tpl, err := LoadTheme(themeDir)
	if err != nil {
		t.Error(err)
	}
	t.Skip("fis todo")
	var wb bytes.Buffer
	tpl.ExecuteTemplate(&wb, "md", "nothing")
	html := wb.String()
	t.Log(html)
	if !strings.Contains(html, "<h1>h1</h1>") {
		t.Error("md parse error")
	}
}

func TestRequireFunc(t *testing.T) {
	themeDir, _ := filepath.Abs("test/output/")
	tpl, err := LoadTheme(themeDir)
	if err != nil {
		t.Error(err)
	}
	t.Skip("fis todo")
	var wb bytes.Buffer
	tpl.ExecuteTemplate(&wb, "require", "nothing")
	html := wb.String()
	t.Log(html)
	if !strings.Contains(html, "<script src=\"test.js\"></script>") {
		t.Error("require js error")
	}
	if !strings.Contains(html, "<link type=\"text/css\" rel=\"stylesheet\" href=\"test.css\">") {
		t.Error("require css error")
	}

}
