package main

import (
	"fmt"
	"image/jpeg"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

type Photo struct {
	Hash *goimagehash.ImageHash
	Path string
}

func main() {
	var photos []Photo

	filepath.Walk("images2", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) == ".jpg" || strings.ToLower(filepath.Ext(path)) == ".jpeg" {
			if h := ahash(path); h != nil {
				photos = append(photos, Photo{
					Hash: h,
					Path: path,
				})
			}
		}
		return nil
	})

	//a b c d e f g
	var groups [][]Photo

	fmt.Println(photos)

	for i := 0; i < len(photos); i++ {
		var group []Photo
		group = append(group, photos[i])

		for j := i + 1; j < len(photos); j++ {
			d := dist(photos[i].Hash, photos[j].Hash)
			fmt.Println(d, photos[i].Path, photos[j].Path)
			if d < 5 {
				group = append(group, photos[j])
				photos = DeleteSlice3(photos, j)
				j--
			}
		}
		groups = append(groups, group)
	}
	for i := range groups {
		fmt.Println("--------------------------------")
		for j := range groups[i] {
			fmt.Println(groups[i][j].Path)
		}
		fmt.Println("--------------------------------")
	}
}

func ahash(path string) *goimagehash.ImageHash {
	file1, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file1.Close()

	img1, err := jpeg.Decode(file1)
	if err != nil {
		return nil
	}

	hash1, err := goimagehash.AverageHash(img1)
	if err != nil {
		return nil
	}
	return hash1
}

func dist(h1, h2 *goimagehash.ImageHash) int {
	distance, err := h1.Distance(h2)
	if err != nil {
		return -1
	}
	return distance
}

// DeleteSlice3 删除指定元素。
func DeleteSlice3(a []Photo, idx int) []Photo {
	j := 0
	for i, v := range a {
		if i != idx {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
