package banana

import (
	"bytes"
	"testing"
)

func TestLoad5xx(t *testing.T) {
	t5xx, name := Load5xx()
	var wb bytes.Buffer
	t5xx.ExecuteTemplate(&wb, name, "nothing")
	if wb.String() != "Template file error" {
		t.Log(wb.String())
		t.Error("tpl is not 5xx")
	}
}
