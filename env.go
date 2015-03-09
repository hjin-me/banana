package banana

import (
	"errors"
	"time"
)

type AppCfg struct {
	Env struct {
		ConfRoot string
		Port     string
		Level    string
		Tpl      string
		Timeout  time.Duration
		Statics  string
		Db       map[string]interface{}
	}
}

var (
	ErrFileNotFound = errors.New("file not found")
	ErrNotFile      = errors.New("path is not file ")
	ErrNotDir       = errors.New("path is not dir")
)
