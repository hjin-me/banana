package banana

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
)

var (
	baseCtx context.Context
)

func prepare() {

	confFilename, _ := filepath.Abs("test/app.yaml")
	cfg := loadCfg(confFilename)
	baseCtx = context.WithValue(context.Background(), "cfg", cfg)
}

func TestContextOutput(t *testing.T) {
	prepare()

	testStr := "hello world"
	testContentType := "text/plain"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCtx, _ := context.WithTimeout(baseCtx, time.Second)

		ctx := WithHttp(rawCtx, w, r, map[string]string{})
		ctx.Output(testStr, testContentType)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	ct, ok := res.Header["Content-Type"]
	if !ok {
		t.Error("Response no Content-Type")
	}
	if ct[0] != testContentType {
		t.Error("Content Type not", testContentType)
	}
	if string(greeting) != testStr {
		t.Error("response error")
	}

	t.Logf("%s\n", greeting)

}

func TestContextJSON(t *testing.T) {
	prepare()

	type TD struct {
		Hello string
		Num   []int
	}
	testData := TD{"world", []int{3, 2, 1}}
	testContentType := "application/json"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCtx, _ := context.WithTimeout(baseCtx, time.Second)

		ctx := WithHttp(rawCtx, w, r, map[string]string{})
		ctx.Json(testData)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	ct, ok := res.Header["Content-Type"]
	if !ok {
		t.Error("Response no Content-Type")
	}
	if ct[0] != testContentType {
		t.Error("Content Type not", testContentType)
	}
	tstr, err := json.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}
	if string(greeting) != string(tstr) {
		t.Error("response error")
	}

	t.Logf("%s\n", greeting)

}

func TestContextTpl5xx(t *testing.T) {
	prepare()

	type TD struct {
		Hello string
		Num   []int
	}
	testData := TD{"world", []int{3, 2, 1}}
	testHtml := "Template file error"
	testContentType := "text/html"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCtx, _ := context.WithTimeout(baseCtx, time.Second)

		ctx := WithHttp(rawCtx, w, r, map[string]string{})
		ctx.Tpl("abc.tpl", testData)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", greeting)
	if res.StatusCode != http.StatusInternalServerError {
		t.Error("status code ", res.StatusCode, "!=", http.StatusInternalServerError)
	}
	ct, ok := res.Header["Content-Type"]
	if !ok {
		t.Error("Response no Content-Type")
	}
	if ct[0] != testContentType {
		t.Error("Content Type not", testContentType)
	}
	if string(greeting) != string(testHtml) {
		t.Error("response error")
	}

}

func TestContextTpl(t *testing.T) {
	prepare()

	type TD struct {
		Hello string
		Num   []int
	}
	testData := TD{"world", []int{3, 2, 1}}
	testContentType := "text/html"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCtx, _ := context.WithTimeout(baseCtx, time.Second)

		ctx := WithHttp(rawCtx, w, r, map[string]string{})
		ctx.Tpl("test:page/demo.html", testData)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", greeting)
	ct, ok := res.Header["Content-Type"]
	if !ok {
		t.Error("Response no Content-Type")
	}
	if ct[0] != testContentType {
		t.Error("Content Type not", testContentType)
	}
	if !strings.Contains(string(greeting), "DOCTYPE html") {
		t.Error("response error", string(greeting))
	}

}
