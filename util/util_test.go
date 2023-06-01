package util

import (
	"fmt"
	"testing"
)

func TestPage(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	pages, start, end := Page(len(s), 3, 2)
	fmt.Println(pages, start, end, s[start:end])
}
