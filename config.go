package banana

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
)

var (
	configCacheMap map[string][]byte = make(map[string][]byte)
	mutex          *sync.Mutex       = &sync.Mutex{}
)

func Config(filename string, data interface{}) (err error) {
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
		// case filename = <-confFilepath(ctx, filename):
		case <-ctx.Done():
			err := errors.New("cant find file [" + filename + "]")
			return err
		}
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		bf, err = ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		configCacheMap[filename] = bf
	}

	err = yaml.Unmarshal(bf, data)
	if err != nil {
		return err
	}
	return
}
