package crontab

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	var time string = "12:32:23"
	exp, err := ParseTime(time)
	if err != nil {
		t.Fatalf("parse time failed, cause: %v", err)
	}
	t.Logf("time: %v to expression: %v", time, exp)

	//
	time = "12:32"
	exp1, err := ParseTime(time)
	if err != nil {
		t.Fatalf("parse time failed, cause: %v", err)
	}
	t.Logf("time: %v to expression: %v", time, exp1)

	time = "@hourly"
	exp2, err := ParseTime(time)
	if err != nil {
		t.Fatalf("parse time failed, cause: %v", err)
	}
	t.Logf("time: %v to expression: %v", time, exp2)

}
