package main

import (
	"math/rand"
	"slices"
	"strconv"
	"testing"
	"time"
)

func generateSliceRandDurations(n int) []time.Duration {
	var result []time.Duration

	for i := 0; i < n; i++ {
		result = append(result, time.Duration(rand.Intn(15))*time.Second)
	}

	return result
}

func startSigs(durs []time.Duration) []<-chan interface{} {
	var result []<-chan interface{}

	for _, dur := range durs {
		result = append(result, sig(dur))
	}

	return result
}

func TestOr(t *testing.T) {
	n := rand.Intn(10)

	for i := 0; i < n; i++ {
		t.Run("TestOr"+strconv.Itoa(i), func(t *testing.T) {
			var durs []time.Duration

			for len(durs) == 0 {
				n := rand.Intn(10)
				durs = generateSliceRandDurations(n)
			}

			maxDur := slices.Max(durs)

			start := time.Now()

			sigs := startSigs(durs)

			<-or(
				sigs...,
			)

			workDur := time.Since(start)
			errDeltaDur := 500 * time.Millisecond

			if workDur-maxDur > errDeltaDur {
				t.Errorf("Result was incorrect, got: %v, want: [%v, %v]", workDur, maxDur, maxDur+errDeltaDur)
			}
		})
	}
}
