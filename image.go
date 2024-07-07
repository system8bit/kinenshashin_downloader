package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func loadImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func concatImageFiles(dstPath string, srcPaths [][]string) error {
	srcImages := make([]image.Image, 0)
	width, height := 0, 0
	for _, path := range srcPaths {
		lastDx := 0
		lastDy := 0
		for _, p := range path {
			img, err := loadImageFromFile(p)
			if err != nil {
				return err
			}
			rct := img.Bounds()
			srcImages = append(srcImages, img)
			lastDx = rct.Dx()
			lastDy += rct.Dy()
		}
		height = max(height, lastDy)
		width += lastDx
	}

	dstImage := image.NewRGBA(image.Rect(0, 0, width, height))
	lenX := len(srcPaths[0])
	offsetX, lastDx := 0, 0
	for i, path := range srcPaths {
		offsetY := 0
		for j, _ := range path {
			img := srcImages[i*lenX+j]
			srcRect := img.Bounds()
			draw.Draw(
				dstImage,
				image.Rect(offsetX, offsetY, offsetX+srcRect.Dx(), offsetY+srcRect.Dy()),
				img,
				image.Point{0, 0},
				draw.Over,
			)
			offsetY += srcRect.Dy()
			lastDx = srcRect.Dx()
		}
		offsetX += lastDx
	}

	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := png.Encode(file, dstImage); err != nil {
		return err
	}
	return nil
}
