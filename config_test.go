package banana

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	var (
		err error
	)
	type TestAppYaml struct {
		Env struct {
			Port int
			Tpl  string `banana:"relative"`
		}
	}
	cfg := TestAppYaml{}
	// pwd log
	_, err = Config("test/app.yaml", &cfg)
	if err != nil {
		t.Error(err)
	}
	if cfg.Env.Port != 8088 {
		t.Error("cfg error", cfg)
	}
	x, _ := os.Getwd()
	if cfg.Env.Tpl != filepath.Join(x, "test/output") {
		t.Error(cfg)
	}

	type TestYaml struct {
		Abs  string `yaml:"abs"`
		Conf string `yaml:"conf"`
	}
	tcfg := TestYaml{}
	// absolute log
	filename, _ := filepath.Abs("test/abs.yaml")
	_, err = Config(filename, &tcfg)
	if err != nil {
		t.Error(err)
	}
	if tcfg.Abs != "hello" {
		t.Error("abs cfg err", tcfg)
	}

	// conf base log
	SetBaseDir(filepath.Dir(filename))
	_, err = Config("conf.yaml", &tcfg)
	if err != nil {
		t.Error(err)
	}
	if tcfg.Conf != "world" {
		t.Error("conf base cfg err", tcfg)
	}

	// not exists log
	_, err = Config("test", &cfg)
	if err == nil {
		t.Error("should cause an error")
	}

}
func TestReflect(t *testing.T) {
	t.Skip()

	type TestAppYaml struct {
		Env struct {
			Port string `banana:"relative"`
		}
		T string `banana:"absolute"`
	}

	data := TestAppYaml{}
	data.Env.Port = "test"
	data.T = "1234"

	scan(reflect.ValueOf(&data))
	t.Log(data)

	data = TestAppYaml{}
	data.Env.Port = "zxc"
	data.T = "123"
	scan(reflect.ValueOf(data))

	t.Log(data)

	t.Error("end")

}
func scan(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	log.Println(v.Kind())
	switch v.Kind() {
	case reflect.Struct:
		n := v.NumField()
		t := v.Type()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.Struct:
				scan(f)
			case reflect.String:
				tf := t.Field(i)
				if s := tf.Tag.Get("banana"); s != "" && f.CanSet() {
					f.SetString("xxxxxx:" + f.String())
				}
			}
		}
	}

}
