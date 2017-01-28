package dither

import (
	"log"
	"os"
	"image"
	"image/png"
	_ "image/jpeg"
	"image/color"
)

func (dither Dither) PrintMono(input image.Image) {
	bounds := input.Bounds()
	img := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
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
			pix -= errors[x][y] * dither.ErrorMultiplier

			var quantError float32
			// Diffuse the error of each calculation to the neighboring pixels
			if pix < 128 {
				quantError = -pix
				pix = 0
			} else {
				quantError = 255 - pix
				pix = 255
			}

			img.SetGray(x, y, color.Gray{Y:uint8(pix)})

			// Diffuse error in two dimension
			ydim := len(dither.Filter) - 1
			xdim := len(dither.Filter[0]) / 2
			for xx := 0; xx < ydim + 1; xx++ {
				for yy := -xdim; yy <= xdim - 1; yy++ {
					if y + yy < 0 || dy <= y + yy || x + xx < 0 || dx <= x + xx {
						continue
					}
					// Adds the error of the previous pixel to the current pixel
					errors[x+xx][y+yy] += quantError * dither.Filter[xx][yy + ydim]
				}
			}
		}
	}
	output, err := os.Create("output/mono/" + dither.Type +".png")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()
	err = png.Encode(output, img)

	if err != nil {
		log.Fatal(err)
	}
}
