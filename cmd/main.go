package main

import (
	"os"
	"flag"
	"fmt"
	"log"
	"time"
	"image"
	"image/color"
	"image/png"
	_ "image/jpeg"
	"github.com/esimov/dithergo"
)

type file struct {
	name string
}

var (
	ditherers 	[]dither.Dither

	// Command line flags
	outputDir	string
	export		string
	grayscale	bool
	treshold	bool
	multiplier	float64
	commands 	flag.FlagSet
)

func (file *file) Open() (image.Image, error) {
	f, err := os.Open(file.name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// Convert image to grayscale
func (file *file) Grayscale(input image.Image, grayscale bool) (*image.Gray, error) {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)

	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			gray.Set(x, y, pixel)
		}
	}
	if grayscale {
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

// Create treshold image
func (file *file) TresholdDithering(input *image.Gray, treshold bool) (*image.Gray, error) {
	var (
		bounds = input.Bounds()
		dithered = image.NewGray(bounds)
		dx = bounds.Dx()
		dy = bounds.Dy()
	)
	if !treshold {
		return nil, nil
	}
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
	output, err := os.Create("output/treshold.png")
	if err != nil {
		return nil, err
	}
	defer output.Close()
	err = png.Encode(output, dithered)

	if err != nil {
		log.Fatal(err)
	}
	return dithered, nil
}

// Function to visualize the rendering progress
func progress(done chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 200)

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
	commands = *flag.NewFlagSet("commands", flag.ExitOnError)
	commands.StringVar(&outputDir, "outputdir", "output", "Directory name, where to save the generated images")
	commands.StringVar(&export, "export", "all", "Generate the color and greyscale dithered images. Options: 'all', 'color', 'mono'")
	commands.BoolVar(&grayscale, "grayscale", true, "Convert image to grayscale")
	commands.BoolVar(&treshold, "treshold", true, "Export treshold image")
	commands.Float64Var(&multiplier, "multiplier", 1.18, "Error multiplier")

	if len(os.Args) <= 1 {
		fmt.Println("Please provide an image, or type --help for the supported command line arguments\n")
		os.Exit(1)
	}

	if (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(`
Usage of commands:
  -export string
    	Generate the color and greyscale dithered images. Options: 'all', 'color', 'mono' (default "all")
  -grayscale
    	Convert image to grayscale (default true)
  -multiplier float
    	Error multiplier (default 1.18)
  -outputdir string
    	Directory name, where to save the generated images (default "output")
  -treshold
    	Export treshold image (default true)
		`)
		os.Exit(1)
	}
	commands.Parse(os.Args[2:])

	// Dithering methods
	ditherers = []dither.Dither{
		dither.Dither{
			"FloydSteinberg",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0 },
					[]float32{ 3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0 },
					[]float32{ 1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0 },
				},
				float32(multiplier),
			},
		},
		dither.Dither{
			"Stucki",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0 },
					[]float32{ 2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0 },
					[]float32{ 1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0 },
				},
				float32(multiplier),
			},
		},
		dither.Dither{
			"Atkinson",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 1.0 / 8.0, 1.0 / 8.0 },
					[]float32{ 1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0.0 },
					[]float32{ 0.0, 1.0 / 8.0, 0.0, 0.0 },
				},
				float32(multiplier),
			},
		},
		dither.Dither{
			"Burkes",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 0.0, 8.0 / 32.0, 4.0 / 32.0 },
					[]float32{ 2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0 },
					[]float32{ 0.0, 0.0, 0.0, 0.0, 0.0 },
				},
				float32(multiplier),
			},
		},
		dither.Dither{
			"Sierra-3",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0 },
					[]float32{ 2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0 },
					[]float32{ 0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0 },
				},
				0.92,
			},
		},
		dither.Dither{
			"Sierra-2",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 0.0, 4.0 / 16.0, 3.0 / 16.0 },
					[]float32{ 1.0 / 16.0, 2.0 / 16.0, 3.0 / 16.0, 2.0 / 16.0, 1.0 / 16.0 },
					[]float32{ 0.0, 0.0, 0.0, 0.0, 0.0 },
				},
				float32(multiplier),
			},
		},
		dither.Dither{
			"Sierra-Lite",
			dither.Settings{
				[][]float32{
					[]float32{ 0.0, 0.0, 2.0 / 4.0 },
					[]float32{ 1.0 / 4.0, 1.0 / 4.0, 0.0 },
					[]float32{ 0.0, 0.0, 0.0 },
				},
				float32(multiplier),
			},
		},
	}

	done := make(chan struct{})
	input := &file{name: string(os.Args[1])}
	img, _ := input.Open()

	fmt.Print("Rendering image...")
	now := time.Now()
	progress(done)

	// Run dither methods
	func(input *file, done chan struct{}) {
		if commands.Parsed() {
			_ = os.Mkdir(outputDir + "/color", os.ModePerm)
			_ = os.Mkdir(outputDir + "/mono", os.ModePerm)
			gray, _ := input.Grayscale(img, grayscale)
			input.TresholdDithering(gray, treshold)

			for _, ditherer := range ditherers {
				switch export {
				case "all":
					ditherer.PrintColor(img)
					ditherer.PrintMono(img)
				case "color":
					ditherer.PrintColor(img)
				case "mono":
					ditherer.PrintMono(img)
				}
			}
			done <- struct{}{}
		}
	}(input, done)

	since := time.Since(now)
	fmt.Println("\nDone✓")
	fmt.Printf("Rendered in: %.2fs\n", since.Seconds())
}