package banana

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/net/context"
)

type MuxContext struct {
	context.Context
}

func (p *MuxContext) Conf() AppCfg {
	cfg, _ := p.Value("cfg").(AppCfg)
	return cfg
}

func (p *MuxContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	statusCode := http.StatusOK
	var err interface{}
	defer func() {
		endTime := time.Now()
		if err == nil {
			err = ""
		}
		log.Printf("%d %0.3f [%s] [%s] [%s]", statusCode, float32(endTime.Sub(startTime))/float32(time.Second), r.URL.Path, r.URL.RawQuery, err)
	}()

	list, exist := routeList[r.Method]
	if !exist {
		statusCode = http.StatusNotFound
		http.NotFound(w, r)
		return
	}

	var (
		ctx     context.Context
		timeout bool = true
		cancel  func()
	)
	if p.Conf().Env.Timeout == 0 {
		ctx, cancel = context.WithCancel(p)
	} else {
		ctx, cancel = context.WithTimeout(p, p.Conf().Env.Timeout*time.Millisecond)
	}
	defer cancel()

	var (
		ruleFound = false
	)
	for _, v := range list {
		res := v.regex.FindStringSubmatch(r.URL.Path)

		params := make(map[string]string)
		for k, v := range v.nameList {
			if len(res) > k+1 {
				params[v] = res[k+1]
			} else {
				params[v] = ""
			}
		}
		if len(res) > 0 {
			ruleFound = true
			go func() {
				defer func() {
					if err = recover(); err != nil {
						statusCode = http.StatusInternalServerError
						w.WriteHeader(statusCode)
					}
					timeout = false
					cancel()
				}()
				e := v.controller(WithHttp(ctx, w, r, params))
				if e != nil {
					panic(e)
				}
			}()
			break
		}
	}
	if !ruleFound {
		timeout = false
		cancel()
	}
	<-ctx.Done()
	switch {
	case timeout:
		err = ctx.Err()
		statusCode = http.StatusGatewayTimeout
		w.WriteHeader(statusCode)
	case !ruleFound:
		statusCode = http.StatusNotFound
		w.WriteHeader(statusCode)
	}
	return
}

func App(args ...string) *MuxContext {
	var (
		ctx *MuxContext
	)
	if len(args) == 0 {
		ctx = initial()
	} else {
		ctx = bootstrap(args[0])
	}

	go func() {
		err := http.ListenAndServe(":"+ctx.Conf().Env.Port, ctx) //设置监听的端口
		if err != nil {
			log.Print(err)
		}
	}()
	return ctx
}

func initial() *MuxContext {
	return bootstrap(flagParams())
}

func bootstrap(confFilename string) *MuxContext {
	runtime.GOMAXPROCS(runtime.NumCPU())
	routeList = make(map[string][]routeInfo)
	cfg := AppCfg{}

	absFilename, err := Config(confFilename, &cfg)
	if err != nil {
		panic(err)
	}
	SetBaseDir(filepath.Dir(absFilename))

	return &MuxContext{context.WithValue(context.Background(), "cfg", cfg)}
}
