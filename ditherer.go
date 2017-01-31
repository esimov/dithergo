package dither

import (
	"os"
	"flag"
	"fmt"
	"log"
	"time"
	"image"
	"errors"
	"image/color"
	"image/png"
	_ "image/jpeg"
)

type file struct {
	name string
}

// Command line flags
var (
	outputDir	string
	export		string
	grayscale	bool
	treshold	bool
	multiplier	float64
	commands 	flag.FlagSet
)

const helper = `
Usage:
go run <image>

Options:
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
`

// Open the input file
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

// Parse the command line inputs and call the dithering methods
func Process(ditherers []Dither)  {
	commands = *flag.NewFlagSet("commands", flag.ExitOnError)
	commands.StringVar(&outputDir, "outputdir", "output", "Directory name, where to save the generated images")
	commands.StringVar(&export, "export", "all", "Generate the color and greyscale dithered images. Options: 'all', 'color', 'mono'")
	commands.BoolVar(&grayscale, "grayscale", true, "Convert image to grayscale")
	commands.BoolVar(&treshold, "treshold", true, "Export treshold image")
	commands.Float64Var(&multiplier, "multiplier", 1.18, "Error multiplier")

	if len(os.Args) <= 1 || (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(errors.New(helper))
		os.Exit(1)
	}

	// Parse flags before to use them
	commands.Parse(os.Args[2:])

	// Channel to signal the completion event
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
				outputColor := ditherer.Color(img, float32(multiplier))
				outputMono  := ditherer.Monochrome(img, float32(multiplier))
				colorExport := outputDir + "/color/"
				monoExport  := outputDir + "/mono/"

				switch export {
				case "all":
					generateOutput(ditherer, outputColor, colorExport)
					generateOutput(ditherer, outputMono, monoExport)
				case "color":
					generateOutput(ditherer, outputColor, colorExport)
				case "mono":
					generateOutput(ditherer, outputMono, monoExport)
				}
			}
			done <- struct{}{}
		}
	}(input, done)

	since := time.Since(now)
	fmt.Println("\nDone✓")
	fmt.Printf("Rendered in: %.2fs\n", since.Seconds())
}

// Output the resulting image
func generateOutput(dither Dither, img image.Image, exportDir string) {
	output, err := os.Create(exportDir + dither.Type +".png")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	err = png.Encode(output, img)
	if err != nil {
		log.Fatal(err)
	}
}