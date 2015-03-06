package banana

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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

func checkDir(base, in string) (string, error) {

	if !filepath.IsAbs(in) {
		in = filepath.Join(base, in)
	}
	fi, err := os.Lstat(in)
	if err != nil {
		log.Println(err, base, in)
		return "", err
	}
	if !fi.IsDir() {
		emsg := fmt.Sprintf("%s: should be directory\n", in)
		log.Printf(emsg)
		return "", errors.New(emsg)
	}

	return in, nil
}
func loadCfg(filename string) (cfg AppCfg) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalln("config file path error", err)
	}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("open config file failed", err)
	}
	defer f.Close()
	bf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("read config file failed", err)
	}
	err = yaml.Unmarshal(bf, &cfg)
	if err != nil {
		log.Fatalln("load config fail", err)
	}
	cfg.Env.ConfRoot = filepath.Dir(filename)
	cfg.Env.Tpl, err = checkDir(cfg.Env.ConfRoot, cfg.Env.Tpl)
	if err != nil {
		log.Fatalln(err)
	}
	return

}
func flagParams() (confFilename string) {
	f := flag.NewFlagSet("params", flag.ExitOnError)
	f.StringVar(&confFilename, "c", "./app.yaml", "server configuration")

	if err := f.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	return
}
