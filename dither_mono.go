package main

import (
	"image"
	"image/color"
)

type MonoDitherer interface {
	Algorithm(image.Image) (*image.Gray, error)
}

type Settings struct {
	filter [][]float32
	errorMultiplier float32
}
type MonoDither struct {
	Type string
	Settings
}

func (dither MonoDither) Algorithm(input image.Gray) (*image.Gray, error) {
	bounds := input.Bounds()
	img := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := img.At(x, y)
			img.Set(x, y, pixel)
		}
	}

	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()

	// Prepopulate multidimensional slice
	errors := make([][]float32, dx)
	for x := 0; x < dx; x++ {
		errors[x] = make([]float32, dy)
		for y := 0; y < dy; y++ {
			errors[x][y] = 0
		}
	}

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			pix := float32(img.GrayAt(x, y).Y)
			pix -= errors[x][y] * dither.errorMultiplier

			var quantError float32
			if pix < 128 {
				quantError = -pix
				pix = 0
			} else {
				quantError = 255 - pix
				pix = 255
			}

			img.SetGray(x, y, color.Gray{Y:uint8(pix)})

			for xx := 0; xx < 3; xx++ {
				for yy := -2; yy <= 2; yy++ {
					if y + yy < 0 || dy <= y + yy || x + xx < 0 || dx <= x + xx {
						continue
					}
					errors[x+xx][y+yy] += quantError * dither.filter[xx][yy + 2]
				}
			}
		}
	}
	return img, nil
}
