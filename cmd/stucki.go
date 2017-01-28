package main

import (
	"image"
	_ "image/jpeg"
	"os"
	"fmt"
	"time"
	"github.com/esimov/dithergo"
)

type file struct {
	name string
}

var ditherers []dither.Dither

func (file *file) Open() (image.Image, error) {
	f, err := os.Open(file.name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
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
			}
		}
	}()
}

func main()  {
	ditherers = []dither.Dither{
		dither.Dither{
			"Stucki",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0 },
					[]float32{ 2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0 },
					[]float32{ 1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0 },
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

	input := &file{name: string(os.Args[1])}
	img, _ := input.Open()
	fmt.Print("Rendering image...")

	progress(done)

	func(input *file, done chan struct{}) {
		_ = os.Mkdir("output/color", os.ModePerm)
		_ = os.Mkdir("output/mono", os.ModePerm)

		for _, ditherer := range ditherers {
			ditherer.PrintColor(img)
			ditherer.PrintMono(img)
		}
		done <-struct{}{}
	}(input, done)

	fmt.Println("\nDone✓")
}