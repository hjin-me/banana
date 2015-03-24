package banana

import (
	"path/filepath"
	"testing"
)

func TestFisLoadMap(t *testing.T) {
	themeDir, _ := filepath.Abs("test/output")
	fisLoadMap(themeDir)
	if v, ok := globalFisMap.Res["test:page/demo.html"]; ok {
		t.Log(v)
		if v.Type != "html" {
			t.Error("type error")
		}
		if !v.IsPage() {
			t.Error("should isPage")
		}
	} else {
		t.Log(globalFisMap)
		t.Error("test:page/demo.html not exists")
	}
	if v, ok := globalFisMap.Res["test:static/lib/mod.js"]; ok {
		t.Log(v)
		if v.Type != "js" {
			t.Error("type error")
		}
		if v.IsPage() {
			t.Error("isPage is false")
		}
	} else {
		t.Log(globalFisMap)
		t.Error("test:static/lib/mod.js not exists")
	}

}
