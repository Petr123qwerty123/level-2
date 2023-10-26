package main

import (
	"math"
	"testing"
	"time"
)

func TestPrintCurrentTime001(t *testing.T) {
	actual, _ := PrintCurrentTime()
	expected := time.Now()

	maxDif := time.Second.Seconds()
	dif := expected.Sub(actual).Seconds()

	if math.Abs(dif) >= maxDif {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}
