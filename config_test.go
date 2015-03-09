package banana

import "testing"

func TestConfig(t *testing.T) {
	var (
		err error
	)
	type TestAppYaml struct {
		Env struct {
			Port int
		}
	}
	cfg := TestAppYaml{}
	err = Config("test/app.yaml", &cfg)
	if err != nil {
		t.Error(err)
	}
	if cfg.Env.Port != 8088 {
		t.Error("cfg error", cfg)
	}

	err = Config("test", &cfg)
	if err == nil {
		t.Error("should cause an error")
	}

}
