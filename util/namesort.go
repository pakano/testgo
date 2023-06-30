package util

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Compare(str1, str2 string) bool {
	a, err := UTF8ToGBK(str1)
	if err != nil {
		return false
	}
	b, err := UTF8ToGBK(str2)
	if err != nil {
		return false
	}
	bLen := len(b)
	for idx, chr := range a {
		if idx > bLen-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

// UTF8ToGBK : transform UTF8 rune into GBK byte array
func UTF8ToGBK(src string) ([]byte, error) {
	GB18030 := simplifiedchinese.All[0]
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), GB18030.NewEncoder()))
}
