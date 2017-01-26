package main

import (
	"image"
	"image/color"
	"image/png"
	_ "image/jpeg"
	"os"
	"fmt"
	"log"
	"time"
)

type ImageDrawer interface {
	Draw([][]float64)
}

type File struct {
	Name string
}

var colorDithers []ColorDither

func (file *File) Open() (image.Image, error) {
	f, err := os.Open(file.Name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

func (file *File) Grayscale(input image.Image, createImageOutput bool) (*image.Gray, error) {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)

	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			gray.Set(x, y, pixel)
		}
	}

	if createImageOutput {
		output, err := os.Create("output/grayscale.png")
		if err != nil {
			return nil, err
		}
		defer output.Close()
		err = png.Encode(output, gray)

		if err != nil {
			log.Fatal(err)
		}
	}

	return gray, nil
}

func (file *File) TresholdDithering(input *image.Gray, createImageOutput bool) (*image.Gray, error) {
	var (
		bounds = input.Bounds()
		dithered = image.NewGray(bounds)
		dx = bounds.Dx()
		dy = bounds.Dy()
	)

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			pixel := input.GrayAt(x, y)
			threshold := func(pixel color.Gray)color.Gray {
				if pixel.Y > 123 {
					return color.Gray{Y:255}
				}
				return color.Gray{Y:0}
			}

			dithered.Set(x, y, threshold(pixel))
		}
	}

	if createImageOutput {
		output, err := os.Create("output/treshold.png")
		if err != nil {
			return nil, err
		}
		defer output.Close()
		err = png.Encode(output, dithered)

		if err != nil {
			log.Fatal(err)
		}
	}

	return dithered, nil
}

func FloydSteinbergDitheringColor (input image.Image, errorMultiplier float32) (*image.RGBA, error){
	bounds := input.Bounds()
	img := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			img.Set(x, y, pixel)
		}
	}
	filter := [][]float32{
		[]float32{ 0.0, 0.0, 0.0, 7.0/48.0, 5.0/48.0},
		[]float32{ 3.0/48.0, 5.0/48.0, 7.0/48.0, 5.0/48.0, 3.0/48.0 },
		[]float32{ 1.0/48.0, 3.0/48.0, 5.0/48.0, 3.0/48.0, 1.0/48.0 },
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
			r -= redErrors[x][y] * errorMultiplier
			g -= greenErrors[x][y] * errorMultiplier
			b -= blueErrors[x][y] * errorMultiplier

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
					redErrors[x+xx][y+yy] 	+= qrr * filter[xx][yy + 2]
					greenErrors[x+xx][y+yy] += qrg * filter[xx][yy + 2]
					blueErrors[x+xx][y+yy] 	+= qrb * filter[xx][yy + 2]
				}
			}
		}
	}

	output, err := os.Create("output/color.png")
	if err != nil {
		return nil, err
	}
	defer output.Close()
	err = png.Encode(output, img)

	if err != nil {
		log.Fatal(err)
	}
	return img, nil
}

func FloydSteinbergDithering (img *image.Gray, errorMultiplier float32) (*image.Gray, error) {
	filter := [][]float32{
		[]float32{ 0.0, 0.0, 0.0, 7.0/48.0, 5.0/48.0},
		[]float32{ 3.0/48.0, 5.0/48.0, 7.0/48.0, 5.0/48.0, 3.0/48.0 },
		[]float32{ 1.0/48.0, 3.0/48.0, 5.0/48.0, 3.0/48.0, 1.0/48.0 },
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
			pix -= errors[x][y] * errorMultiplier

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
					errors[x+xx][y+yy] += quantError * filter[xx][yy + 2]
				}
			}
		}
	}

	output, err := os.Create("output/output.png")
	if err != nil {
		return nil, err
	}
	defer output.Close()
	err = png.Encode(output, img)

	if err != nil {
		log.Fatal(err)
	}
	return img, nil
}

func progress(done chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 100)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
			case <-done:
				ticker.Stop()
				fmt.Print("Done!")
			}
		}
	}()
}

func main()  {
	colorDithers = []ColorDither{
		ColorDither{
			"FloydSteinberg",
			Settings{
				[][]float32{
					[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
					[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0 },
					[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0 },
				},
				0.92,
			},
		},
		ColorDither{
			"Stucki",
			Settings{
				[][]float32{
					[]float32{0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0},
					[]float32{2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0 },
					[]float32{1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0 },
				},
				0.92,
			},
		},
	}

	if len(os.Args) < 2 || (len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h")) {
		fmt.Println("Usage :  Command --file name", os.Args)
		os.Exit(1)
	}

	done := make(chan struct{})

	input := &File{Name: string(os.Args[1])}
	img, _ := input.Open()
	fmt.Print("Rendering image...")

	progress(done)

	func(input *File, done chan struct{}) {
		_ = os.Mkdir("output", os.ModePerm)
		grayscale, _ := input.Grayscale(img, false)
		input.TresholdDithering(grayscale, true)
		//FloydSteinbergDithering(grayscale, 0.78)
		//FloydSteinbergDitheringColor(img, 0.92)

		for _, colorDither := range colorDithers{
			dithers := []ColorDitherer{colorDither}
			for _, dither := range dithers {
				switch dither.(type) {
				case ColorDither:
					fmt.Println(colorDither.Type)
					result, err := dither.Algorithm(img)
					if err != nil {
						log.Fatal(err)
					}
					output, err := os.Create("output/" + colorDither.Type +".png")
					defer output.Close()
					err = png.Encode(output, result)

					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
		done <-struct{}{}
	}(input, done)

	fmt.Println("\nDone✓")
}