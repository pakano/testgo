package reverse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	type TeseCase struct {
		Input    string
		Expected string
	}
	testcases := []TeseCase{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"earth", "htrae"},
	}

	for _, testcase := range testcases {
		actual := Reverse(testcase.Input)
		assert.Equal(t, testcase.Expected, actual)
	}
}

func FuzzReverse(f *testing.F) {
	var seeds = []string{"hello", "world", "earth"}
	for i := range seeds {
		f.Add(seeds[i])
	}
	f.Fuzz(func(t *testing.T, input string) {
		if input == "hello" || input == "world" || input == "earth" {
			return
		}
		str1 := Reverse(input)
		str2 := Reverse(str1)
		if strings.EqualFold(input, str2) {
			t.Errorf("reverse failed! input: %s", input)
		}
	})
}
