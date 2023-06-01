package util

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestXxx(t *testing.T) {
	for {
		num := rand.Intn(256)
		x := GetRanStr(num)
		y := Encode([]byte(x))
		z, err := Decoded(y)
		if err != nil {
			t.Error(err)
		}
		if x != string(z) {
			t.Error(x)
		}
	}
}

func Test002(t *testing.T) {
	x := "8ca88cb62d551445ff94e90aead002d9ced182dca1b6465b0427c828534b01430ibd5qfjh3m3utnu"
	z, err := Decoded(x)
	fmt.Println(string(z), err)
}
