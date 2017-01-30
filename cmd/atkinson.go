package main

import (
	"os"
	"flag"
	"fmt"
	"time"
	"image"
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

	// Atkinson method
	ditherers = []dither.Dither{
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