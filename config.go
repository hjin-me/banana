package banana

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

var (
	configCacheMap map[string][]byte = make(map[string][]byte)
	mutex          *sync.Mutex       = &sync.Mutex{}
	baseConfDir    string
)

func configScan(v reflect.Value, base string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		n := v.NumField()
		t := v.Type()
		for i := 0; i < n; i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.Struct:
				configScan(f, base)
			case reflect.String:
				tf := t.Field(i)
				if s := tf.Tag.Get("banana"); s == "relative" && f.CanSet() {
					f.SetString(filepath.Join(base, f.String()))
				}
			}
		}
	}
}
func configUnmarshal(bf []byte, data interface{}, filename string) (err error) {
	err = yaml.Unmarshal(bf, data)
	if err != nil {
		return
	}
	rv := reflect.ValueOf(data)
	configScan(rv, filepath.Dir(filename))
	return
}

func setBaseDir(dir string) {
	baseConfDir = dir
}

func Config(filename string, data interface{}) (absFilename string, err error) {
	bf, ok := configCacheMap[filename]
	if !ok {
		mutex.Lock()
		defer mutex.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		pwd, _ := os.Getwd()

		select {
		case filename = <-absFilepath(ctx, filename):
		case filename = <-relativeFilepath(ctx, pwd, filename):
		case filename = <-relativeFilepath(ctx, baseConfDir, filename):
		case <-ctx.Done():
			err := errors.New("cant find file [" + filename + "]")
			return "", err
		}
		f, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer f.Close()
		bf, err = ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		configCacheMap[filename] = bf
		absFilename = filename
	}

	err = configUnmarshal(bf, data, filename)
	if err != nil {
		return "", err
	}
	return
}
