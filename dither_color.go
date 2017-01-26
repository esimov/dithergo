package main

import (
	"image"
	"image/color"
)

type ColorDitherer interface {
	Algorithm(image.Image) (*image.RGBA, error)
}

type Settings struct {
	filter [][]float32
	errorMultiplier float32
}
type ColorDither struct {
	Type string
	Settings
}

func (dither ColorDither) Algorithm(input image.Image) (*image.RGBA, error) {
	bounds := input.Bounds()
	img := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			img.Set(x, y, pixel)
		}
	}

	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()

	// Prepopulate multidimensional slices
	redErrors   := make([][]float32, dx)
	greenErrors := make([][]float32, dx)
	blueErrors  := make([][]float32, dx)
	for x := 0; x < dx; x++ {
		redErrors[x]	= make([]float32, dy)
		greenErrors[x]	= make([]float32, dy)
		blueErrors[x]	= make([]float32, dy)
		for y := 0; y < dy; y++ {
			redErrors[x][y]   = 0
			greenErrors[x][y] = 0
			blueErrors[x][y]  = 0
		}
	}

	var qrr, qrg, qrb float32
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			r32,g32,b32,a := img.At(x, y).RGBA()
			r, g, b := float32(uint8(r32)), float32(uint8(g32)), float32(uint8(b32))
			r -= redErrors[x][y] * dither.errorMultiplier
			g -= greenErrors[x][y] * dither.errorMultiplier
			b -= blueErrors[x][y] * dither.errorMultiplier

			if r < 128 {
				qrr = -r
				r = 0
			} else {
				qrr = 255 - r
				r = 255
			}
			if g < 128 {
				qrg = -g
				g = 0
			} else {
				qrg = 255 - g
				g = 255
			}
			if b < 128 {
				qrb = -b
				b = 0
			} else {
				qrb = 255 - b
				b = 255
			}
			img.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})

			for xx := 0; xx < 3; xx++ {
				for yy := -2; yy <= 2; yy++ {
					if y + yy < 0 || dy <= y + yy || x + xx < 0 || dx <= x + xx {
						continue
					}
					redErrors[x+xx][y+yy] 	+= qrr * dither.filter[xx][yy + 2]
					greenErrors[x+xx][y+yy] += qrg * dither.filter[xx][yy + 2]
					blueErrors[x+xx][y+yy] 	+= qrb * dither.filter[xx][yy + 2]
				}
			}
		}
	}
	return img, nil
}
