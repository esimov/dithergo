package dither

import (
	"log"
	"os"
	"image"
	"image/png"
	_ "image/jpeg"
	"image/color"
)

type Settings struct {
	Filter [][]float32
	ErrorMultiplier float32
}

type Dither struct {
	Type string
	Settings
}

func (dither Dither) PrintColor(input image.Image) {
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
			r -= redErrors[x][y] * dither.ErrorMultiplier
			g -= greenErrors[x][y] * dither.ErrorMultiplier
			b -= blueErrors[x][y] * dither.ErrorMultiplier

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

			ydim := len(dither.Filter) - 1
			xdim := len(dither.Filter[0]) / 2

			for xx := 0; xx < ydim + 1; xx++ {
				for yy := -xdim; yy <= xdim - 1; yy++ {
					if y + yy < 0 || dy <= y + yy || x + xx < 0 || dx <= x + xx {
						continue
					}
					redErrors[x+xx][y+yy] 	+= qrr * dither.Filter[xx][yy + ydim]
					greenErrors[x+xx][y+yy] += qrg * dither.Filter[xx][yy + ydim]
					blueErrors[x+xx][y+yy] 	+= qrb * dither.Filter[xx][yy + ydim]
				}
			}
		}
	}
	output, err := os.Create("output/color/" + dither.Type +".png")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()
	err = png.Encode(output, img)

	if err != nil {
		log.Fatal(err)
	}
}
