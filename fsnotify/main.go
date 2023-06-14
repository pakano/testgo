package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/utils/inotify"
)

var WatchFlag = inotify.InDelete | inotify.InCloseWrite | inotify.InMove | inotify.InCreate | inotify.InAttrib
var w *inotify.Watcher

func main() {
	var err error
	w, err = inotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	w.AddWatch("/sata", WatchFlag)
	for {
		select {
		case ev := <-w.Event:
			fmt.Println(ev.Name)
			ProcessEvent(ev)
		case err = <-w.Error:
			log.Println("error:", err)
		}
	}
}

func ProcessEvent(ev *inotify.Event) {
	isdir := 0
	if ev.Mask&inotify.InIsdir == inotify.InIsdir {
		isdir = 1
	}

	switch {
	case ev.Mask&inotify.InCreate == inotify.InCreate:
		HandleCreateEvent(ev.Name, isdir)
	}
}

func HandleCreateEvent(path string, isdir int) {
	fileInfo, e := os.Lstat(path)
	if e != nil {
		log.Println("err:", e, path)
		return
	}
	if isdir == 1 {
		RecursiveAddPath(path, fileInfo)
	}
}

func RecursiveAddPath(path string, fileInfo os.FileInfo) {
	if fileInfo.IsDir() {
		time.Sleep(time.Second)
		if err := w.AddWatch(path, WatchFlag); err != nil {
			log.Println("err:", err)
		}
		fileInfos, err := ioutil.ReadDir(path)
		if err != nil {
			log.Println("err:", err)
			return
		}
		for _, fileInfo := range fileInfos {
			path := filepath.Join(path, fileInfo.Name())
			fileInfo, e := os.Lstat(path)
			if e != nil {
				log.Println("err:", e, path)
				continue
			}
			RecursiveAddPath(path, fileInfo)
		}
	}
}
