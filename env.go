package banana

import "time"

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
