package banana

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestCustomMuxHttpHandle(t *testing.T) {
	testStr := `{"x":"abc"}`
	testContentType := "application/json"

	confFilename, _ := filepath.Abs("test/app.yaml")
	cm := bootstrap(confFilename)
	Get("/test/:x", func(ctx Context) error {
		ctx.Json(ctx.Params())
		return nil
	})

	ts := httptest.NewServer(http.HandlerFunc(cm.ServeHTTP))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test/abc")
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
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

func TestTimeout(t *testing.T) {

	confFilename, _ := filepath.Abs("test/app.yaml")
	cm := bootstrap(confFilename)
	Get("/test/:x", func(ctx Context) error {
		time.Sleep(10 * time.Second)
		ctx.Json(ctx.Params())
		return nil
	})

	ts := httptest.NewServer(http.HandlerFunc(cm.ServeHTTP))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test/abc")
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusGatewayTimeout {
		t.Error("status code ", res.StatusCode, "!=", http.StatusGatewayTimeout)
	}

}

func TestNotFound(t *testing.T) {

	confFilename, _ := filepath.Abs("test/app.yaml")
	cm := bootstrap(confFilename)

	Get("/exists", func(ctx Context) error {
		ctx.Json(ctx.Params())
		return nil
	})

	ts := httptest.NewServer(http.HandlerFunc(cm.ServeHTTP))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/notExists")
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Error("status code ", res.StatusCode, "!=", http.StatusNotFound)
	}
}
