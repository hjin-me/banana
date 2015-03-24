package banana

import (
	"path/filepath"
	"testing"
)

func TestFisLoadMap(t *testing.T) {
	themeDir, _ := filepath.Abs("test/output/config")
	fisLoadMap(themeDir)
	if v, ok := globalFisMap.Res["test:page/demo.html"]; ok {
		t.Log(v)
		if v.Type != "html" {
			t.Error("type error")
		}
	} else {
		t.Log(globalFisMap)
		t.Error("test:page/demo.html not exists")
	}
}
