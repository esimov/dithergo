# dithergo

<strong>Dithergo</strong> is a simple Go library implementing various dithering algorithms to produce halftone images. It supports color and monochrome image outputs.

It's main purpose is mostly to explore the Go language as a medium for artistic exploration, but it can be used for more pragmatic purposes too, like to prepare a color image for output on non-color printer. 

Currently the library includes the following dithering algorithms: <strong><i>Floyd Steinberg, Atkinson, Burkes, Stucki, Sierra-2, Sierra-3, Sierra-Lite</i></strong>. All of these algorithms have something in common: they diffuse the error in two dimensions, but they always push the error forward, never backward.

We can represent this with the following diagram:

             X   7   5 
     3   5   7   5   3
     1   3   5   3   1

           (1/48)

where `X` represent the current pixel processed. The fraction at the bottom represents the divisor for the error. Above is the  the `Floyd-Steinberg` dithering algorithm, which is transposed in Go like this:

```go
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
}
```

You can plug in any dithering algorithm, so the library can be further extended.

### Installation

`go get -u github.com/esimov/dithergo`

### Running

Type `go run cmd/main.go --help` to check all the supported commands. The library supports the following commands:

```
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
```
You can run all the dithering algorithms at once, or you can run a specific one from the `cmd` directory.  

### Examples:
|  Input  |
|:--:|
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/input/david.jpg" height="250">|

The below images are generated with the default options using Michelangelo's David statue as sample image.

|  Color  | Monochrome   |
|:--:|:--:|
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/Atkinson.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/Atkinson.png" height="250"> |
Atkinson | Atkinson |
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/Burkes.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/Burkes.png" height="250"> |
Burkes | Burkes |
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/FloydSteinberg.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/FloydSteinberg.png" height="250"> |
Floyd-Steinberg | Floyd-Steinberg | 
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/Sierra-2.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/Sierra-2.png" height="250"> | 
Sierra-2 | Sierra-2 | 
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/Sierra-3.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/Sierra-3.png" height="250"> |
Sierra-3 | Sierra-3 | 
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/Sierra-Lite.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/Sierra-Lite.png" height="250"> |
Sierra-Lite | Sierra-Lite | 
|<img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/color/Stucki.png" height="250"> | <img src="https://raw.githubusercontent.com/esimov/dithergo/master/output/mono/Stucki.png" height="250"> |
Stucki | Stucki | 

## License

This software is distributed under the MIT license found in the LICENSE file.
