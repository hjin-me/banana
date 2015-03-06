package banana

import "testing"

func TestParseRule(t *testing.T) {
	_, s, err := parseRule("/abc/:d")
	if err != nil {
		t.Fatal(err)
	}
	targetParams := []string{"d"}
	if len(s) != len(targetParams) {
		t.Error("length not equal", s, targetParams)
	}
	for k, v := range s {
		if v != targetParams[k] {
			t.Error("parse failed")
		}
	}

}
