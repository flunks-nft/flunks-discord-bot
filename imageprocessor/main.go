package imageprocessor

import (
	"fmt"
	"image"
	"image/draw"
	"os"
)

func Concatenate(image1, image2 string) (image.Image, error) {
	// Load the two PNG images
	img1, err := loadImage(image1)
	if err != nil {
		return nil, fmt.Errorf("Error loading image 1: %w", err)
	}

	img2, err := loadImage(image2)
	if err != nil {
		return nil, fmt.Errorf("Error loading image 2: %w", err)
	}

	// Get the dimensions of the two images
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	// Calculate the dimensions of the concatenated image
	width := bounds1.Dx() + bounds2.Dx()
	height := bounds1.Dy()

	// Create a new blank image to hold the concatenated image
	concatenatedImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw the first image onto the concatenated image
	draw.Draw(concatenatedImage, bounds1, img1, image.Point{0, 0}, draw.Src)

	// Draw the second image onto the concatenated image at the right position
	draw.Draw(concatenatedImage, bounds2.Add(image.Point{bounds1.Dx(), 0}), img2, image.Point{0, 0}, draw.Src)

	return concatenatedImage, nil
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
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
