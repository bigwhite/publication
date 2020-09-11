package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	// 支持png, jpeg, gif
	width, height, err := imageSize(os.Args[1])
	if err != nil {
		fmt.Println("get image size error:", err)
		return
	}
	fmt.Printf("image size: [%d, %d]\n", width, height)
}

func imageSize(imageFile string) (int, int, error) {
	f, _ := os.Open(imageFile)
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return 0, 0, err
	}

	b := img.Bounds()
	return b.Max.X, b.Max.Y, nil
}
