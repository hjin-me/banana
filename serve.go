package banana

import (
	"log"
	"net/http"
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
	list, exist := routeList[r.Method]
	if !exist {
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
					if err := recover(); err != nil {
						log.Println(err)
						w.WriteHeader(http.StatusInternalServerError)
					}
					timeout = false
					cancel()
				}()
				v.controller(WithHttp(ctx, w, r, params))
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
		w.WriteHeader(http.StatusGatewayTimeout)
	case !ruleFound:
		w.WriteHeader(http.StatusNotFound)
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

	cfg := loadCfg(confFilename)

	return &MuxContext{context.WithValue(context.Background(), "cfg", cfg)}
}
