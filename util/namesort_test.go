package util

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

func GetMixStr(n int) string {
	a := make([]rune, n)
	for i := range a {
		kind := rand.Intn(3)
		switch kind {
		case 0:
			a[i] = rune(RandInt(48, 58))
		case 1:
			a[i] = rune(RandInt(65, 91))
		case 2:
			a[i] = rune(RandInt(19968, 40869))
		}
	}
	return string(a)
}

// RandInt [min,max)
func RandInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}

// go test -fuzz=Fuzz -run=FuzzReverse
func FuzzReverse(f *testing.F) {
	var seeds = []int{1, 2, 3, 4, 5, 6}
	for i := range seeds {
		f.Add(seeds[i])
	}
	f.Fuzz(func(t *testing.T, input int) {
		strs := make([]string, 20)
		for i := 0; i < len(strs); i++ {
			strs[i] = GetMixStr(32)
		}

		sort.Slice(strs, func(i, j int) bool {
			return Compare(strs[i], strs[j])
		})

		last := 0
		for i := range strs {
			if i == 0 {
				continue
			}
			if Compare(strs[i], strs[last]) {
				t.Errorf("error\n")
			}
			last = i
		}
	})
}
