package banana

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"golang.org/x/net/context"
)

var (
	extMap = make(map[string]string)
)

type FisRes struct {
	URI    string                 `json:"uri""`
	Type   string                 `json:"type"`
	Extras map[string]interface{} `json:"extras"`
}

type FisMap struct {
	Res map[string]FisRes `json:"res"`
	Pkg struct{}          `json:"pkg"`
}

func fisScanMapDir(ctx context.Context, dir string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		fileList, _ := ioutil.ReadDir(dir)
		for _, v := range fileList {
			ext := filepath.Ext(v.Name())
			if ext != ".json" {
				continue
			}
			fp := filepath.Join(dir, v.Name())
			if !v.IsDir() {
				select {
				case <-ctx.Done():
					return
				case ch <- fp:
				}
			}
		}
	}()
	return ch
}
func fisParseMap(ctx context.Context, ch <-chan string) <-chan FisMap {
	fmCh := make(chan FisMap)
	go func() {
		defer close(fmCh)
		for filename := range ch {
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println(err)
				continue
			}

			fm := FisMap{}
			err = json.Unmarshal(b, &fm)
			if err != nil {
				log.Println(err)
				continue
			}
			select {
			case <-ctx.Done():
				return
			case fmCh <- fm:
			}

		}
	}()
	return fmCh
}
func fisMergeMap(ctx context.Context, chs ...<-chan FisMap) FisMap {
	gFM := FisMap{make(map[string]FisRes), struct{}{}}
	var wg sync.WaitGroup
	wg.Add(len(chs))

	handle := func(ch <-chan FisMap) {
		for fm := range ch {
			for k, v := range fm.Res {
				gFM.Res[k] = v
			}
		}
		wg.Done()
	}

	for _, ch := range chs {
		go handle(ch)
	}

	wg.Wait()
	return gFM
}

var globalFisMap FisMap

func fisLoadMap(dir string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := fisScanMapDir(ctx, dir)
	globalFisMap = fisMergeMap(ctx, fisParseMap(ctx, in), fisParseMap(ctx, in))
}

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
