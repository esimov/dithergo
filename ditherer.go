package dither

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"time"
)

type file struct {
	name string
}

// Command line flags
var (
	outputDir  string
	colorType  string
	treshold   bool
	multiplier float64
	cmd        flag.FlagSet
)

// Open reads the source image
func (file *file) Open() (image.Image, error) {
	f, err := os.Open(file.name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// Grayscale converts an image to grayscale mode
func (file *file) Grayscale(input image.Image) (*image.Gray, error) {
	bounds := input.Bounds()
	gray := image.NewGray(bounds)

	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			pixel := input.At(x, y)
			gray.Set(x, y, pixel)
		}
	}
	return gray, nil
}

// tresholdDithering creates a tresholded image
func (file *file) tresholdDithering(input *image.Gray) (*image.Gray, error) {
	var (
		bounds   = input.Bounds()
		dithered = image.NewGray(bounds)
		dx       = bounds.Dx()
		dy       = bounds.Dy()
	)

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			pixel := input.GrayAt(x, y)
			threshold := func(pixel color.Gray) color.Gray {
				if pixel.Y > 123 {
					return color.Gray{Y: 255}
				}
				return color.Gray{Y: 0}
			}

			dithered.Set(x, y, threshold(pixel))
		}
	}
	output, err := os.Create(outputDir + "/treshold.png")
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

// Process parses the command line inputs and calls the defined dithering method
func Process(ditherers []Dither) {
	cmd = *flag.NewFlagSet("commands", flag.ExitOnError)
	cmd.StringVar(&outputDir, "o", "", "Output folder")
	cmd.StringVar(&colorType, "e", "all", "Generates & exports the color and greyscale mode halftone images. Options: 'all', 'color', 'mono'")
	cmd.BoolVar(&treshold, "t", true, "Option to export the tresholded image")
	cmd.Float64Var(&multiplier, "em", 1.18, "Error multiplier")

	cmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <image>\n", os.Args[0])
		cmd.PrintDefaults()
	}

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <image>\n", os.Args[0])
		os.Exit(0)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stderr, "Usage: %s <image>\n", os.Args[0])
		cmd.PrintDefaults()
		os.Exit(0)
	}

	// Parse flags before to use them
	cmd.Parse(os.Args[2:])

	if len(cmd.Args()) > 0 {
		log.Fatal("Use an image file as the first argument!")
	}

	if len(outputDir) == 0 {
		log.Fatal("Please specify an output directory!")
	}

	input := &file{name: string(os.Args[1])}
	img, err := input.Open()
	if err != nil {
		log.Fatalf("Input error: %v", err)
	}

	// Channel to signal the completion event
	done := make(chan struct{})

	fmt.Print("Rendering image...")
	now := time.Now()
	progress(done)

	// Run the ditherer method
	func(input *file, done chan struct{}) {
		if cmd.Parsed() {
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				os.Mkdir(outputDir, os.ModePerm)
			}
			_ = os.Mkdir(outputDir+"/color", os.ModePerm)
			_ = os.Mkdir(outputDir+"/mono", os.ModePerm)

			if treshold {
				gray, _ := input.Grayscale(img)
				input.tresholdDithering(gray)
			}

			for _, ditherer := range ditherers {
				dc := ditherer.Color(img, float32(multiplier))
				dg := ditherer.Monochrome(img, float32(multiplier))
				cex := outputDir + "/color/"
				gex := outputDir + "/mono/"

				switch colorType {
				case "all":
					generateOutput(ditherer, dc, cex)
					generateOutput(ditherer, dg, gex)
				case "color":
					generateOutput(ditherer, dc, cex)
				case "mono":
					generateOutput(ditherer, dg, gex)
				}
			}
			done <- struct{}{}
		}
	}(input, done)

	since := time.Since(now)
	fmt.Println("\nDone✓")
	fmt.Printf("Rendered in: %.2fs\n", since.Seconds())
}

// generateOutput render the generated image
func generateOutput(dither Dither, img image.Image, exportDir string) {
	output, err := os.Create(exportDir + dither.Type + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	err = png.Encode(output, img)
	if err != nil {
		log.Fatal(err)
	}
}

// progress visualize the rendering progress
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
