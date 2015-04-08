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

type FisRes struct {
	URI    string                 `json:"uri""`
	Type   string                 `json:"type"`
	Extras map[string]interface{} `json:"extras"`
}

type FisMap struct {
	Res map[string]FisRes `json:"res"`
	Pkg struct{}          `json:"pkg"`
}

func (fr FisRes) IsPage() bool {
	b, ok := fr.Extras["isPage"]
	if !ok {
		return false
	}
	t, ok := b.(bool)
	if !ok {
		return false
	}
	return t
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
func fisMergeMap(ctx context.Context, dir string, chs ...<-chan FisMap) FisMap {
	gFM := FisMap{make(map[string]FisRes), struct{}{}}
	var wg sync.WaitGroup
	wg.Add(len(chs))

	handle := func(ch <-chan FisMap) {
		for fm := range ch {
			for k, v := range fm.Res {
				if v.IsPage() {
					v.URI = filepath.Join(dir, v.URI)
				}
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

var (
	globalFisMap FisMap
	mapOnce      sync.Once
)

func fisLoadMap(dir string) FisMap {
	mapOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		in := fisScanMapDir(ctx, filepath.Join(dir, "config"))
		globalFisMap = fisMergeMap(ctx, filepath.Join(dir, "template"), fisParseMap(ctx, in), fisParseMap(ctx, in))
	})
	return globalFisMap
}

func fisLink(args ...interface{}) template.HTML {

	m, ok := args[0].(string)
	if !ok {
		panic("args[0] is not string")
	}

	v, ok := globalFisMap.Res[m]
	if !ok {
		panic("resource not found [" + m + "]")
	}
	s := ""

	switch v.Type {
	case "js":
		s = fmt.Sprintf("<script src=\"%s\"></script>", v.URI)
	case "css":
		s = fmt.Sprintf("<link type=\"text/css\" rel=\"stylesheet\" href=\"%s\">", v.URI)
	}

	return template.HTML(s)
}

func fisURI(args ...interface{}) template.HTML {
	m, ok := args[0].(string)
	if !ok {
		panic("args[0] is not string")
	}

	v, ok := globalFisMap.Res[m]
	if !ok {
		panic("resource not found [" + m + "]")
	}

	return template.HTML(v.URI)
}

func TplExists(base string) bool {
	_, ok := globalFisMap.Res[base]
	return ok
}
