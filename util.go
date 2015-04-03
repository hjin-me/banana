package banana

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
)

func absFilepath(ctx context.Context, filename string) <-chan string {
	ch := make(chan string)
	go func() {
		var (
			err error
		)
		defer func() {
			if err != nil {
				<-ctx.Done()
				close(ch)
				return
			}
			select {
			case <-ctx.Done():
				close(ch)
			default:
				ch <- filename
				close(ch)
			}
		}()
		if !filepath.IsAbs(filename) {
			err = errors.New("filename is not abs")
			return
		}
		fi, err := os.Lstat(filename)
		if err != nil {
			return
		}
		if fi.IsDir() {
			err = ErrNotFile
			return
		}
	}()
	return ch
}

func relativeFilepath(ctx context.Context, base, filename string) <-chan string {
	ch := make(chan string)
	go func() {
		var (
			err error
		)
		defer func() {
			if err != nil {
				<-ctx.Done()
				close(ch)
				return
			}
			select {
			case <-ctx.Done():
				close(ch)
			default:
				ch <- filename
				close(ch)
			}
		}()
		// 相对路径检查
		filename = filepath.Join(base, filename)
		fi, err := os.Lstat(filename)
		if err != nil {
			return
		}
		if fi.IsDir() {
			err = ErrNotFile
			return
		}
	}()
	return ch
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

/*
func loadCfg(filename string) (cfg AppCfg) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		log.Fatalln("config file path error", err)
	}
	_, err = Config(filename, &cfg)
	if err != nil {
		log.Fatalln("load config fail", err)
	}
	return
}
*/

func flagParams() (confFilename string) {
	f := flag.NewFlagSet("params", flag.ExitOnError)
	f.StringVar(&confFilename, "c", "./app.yaml", "server configuration")

	if err := f.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	return
}
