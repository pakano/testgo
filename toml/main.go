package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type songInfo struct {
	Name     string
	Duration int
}

type mysqlInfo struct {
	Host     string
	PassWord string
}

type config struct {
	Bc     string
	Names  []int
	Names2 string
	Song   songInfo
	Mysql  mysqlInfo
}

func test_toml() {
	var cg config
	var cpath string = "./test.toml"
	if _, err := toml.DecodeFile(cpath, &cg); err != nil {
		log.Fatal(err)
	}

	log.Println(cg.Bc, cg.Song.Name, cg.Mysql.PassWord, cg.Names, cg.Names2)
}

func main() {
	test_toml()
}
