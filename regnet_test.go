package regnet_test

import (
	"github.com/anaray/regnet"
	"testing"
)

const (
	blockIdent string = "REGNET_BLOCK"
	blockKey   string = "REGNET_KEY"
)

func TestNewRegnet(t *testing.T) {
	r, _ := regnet.New()
	if r == nil {
		t.Errorf("regnet initialization failed")
	}

	value, present := r.GetPattern(blockIdent)
	if present == false {
		t.Errorf("expected pattern %s but received nil", value)
	}

	value, present = r.GetPattern(blockKey)
	if present == false {
		t.Errorf("expected pattern %s but received nil", value)
	}

	key := "%{DAY}"
	patternStr := "(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)"
	r.AddPattern(key, patternStr)
	pattern, _ := r.GetPattern(key)

	if pattern.Compiled.String() != "(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)" {
		t.Errorf("expected pattern %s but received %s", patternStr, pattern.Compiled.String())
	}
}
