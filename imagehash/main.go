package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/corona10/goimagehash"
)

func main() {
	file1, _ := os.Open("images/1.jpg")
	file2, _ := os.Open("images/2.jpg")
	defer file1.Close()
	defer file2.Close()

	img1, _ := jpeg.Decode(file1)
	img2, _ := jpeg.Decode(file2)
	hash1, _ := goimagehash.AverageHash(img1)
	hash2, _ := goimagehash.AverageHash(img2)
	distance, _ := hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)

	hash1, _ = goimagehash.DifferenceHash(img1)
	hash2, _ = goimagehash.DifferenceHash(img2)
	distance, _ = hash1.Distance(hash2)
	fmt.Printf("Distance between images: %v\n", distance)
	width, height := 8, 8
	hash3, _ := goimagehash.ExtAverageHash(img1, width, height)
	hash4, _ := goimagehash.ExtAverageHash(img2, width, height)
	distance, _ = hash3.Distance(hash4)
	fmt.Printf("Distance between images: %v\n", distance)

	fmt.Printf("hash3 bit size: %v\n", hash3.Bits())
	fmt.Printf("hash4 bit size: %v\n", hash4.Bits())

	// var b bytes.Buffer
	// foo := bufio.NewWriter(&b)
	// _ = hash4.Dump(foo)
	// foo.Flush()
	// bar := bufio.NewReader(&b)
	// hash5, _ := goimagehash.LoadExtImageHash(bar)
}
