package banana

import (
	"errors"
	"time"
)

type AppCfg struct {
	Env struct {
		ConfRoot string `banana:"relative"`
		Port     string
		Level    string
		Tpl      string `banana:"relative"`
		Timeout  time.Duration
		Statics  string `banana:"relative"`
		Db       map[string]interface{}
	}
}

var (
	ErrFileNotFound = errors.New("file not found")
	ErrNotFile      = errors.New("path is not file ")
	ErrNotDir       = errors.New("path is not dir")
)
