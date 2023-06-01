package main

import "regexp"

var reg *regexp.Regexp

func init() {
	var err error
	reg, err = regexp.Compile("^\\/(\\w+\\/?)+$")
	if err != nil {
		panic(err)
	}
}

func main() {

}
