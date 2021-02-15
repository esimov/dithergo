# dithergo
[![CI](https://github.com/esimov/dithergo/workflows/CI/badge.svg)](https://github.com/esimov/dithergo/actions)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go)](https://pkg.go.dev/github.com/esimov/dithergo)
[![Go Reference](https://pkg.go.dev/badge/github.com/esimov/dithergo.svg)](https://pkg.go.dev/github.com/esimov/dithergo)
[![license](https://img.shields.io/github/license/esimov/dithergo)](./LICENSE)

<strong>Dithergo</strong> is a simple Go library implementing various [dithering](https://en.wikipedia.org/wiki/Dither) algorithm to produce halftone images. It supports color and monochrome image outputs.

The library implements the following dithering methods: ***Floyd Steinberg, Atkinson, Burkes, Stucki, Sierra-2, Sierra-3, Sierra-Lite***. All of these algorithms have something in common: they diffuse the error in two dimensions, but they always push the error forward, never backward.

We can represent this with the following diagram:

             X   7   5 
     3   5   7   5   3
     1   3   5   3   1

           (1/48)

where `X` represent the current pixel processed. The fraction at the bottom represents the divisor for the error. Above is the  the `Floyd-Steinberg` dithering algorithm which can be transposed into the following Go code:

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

`$ go get -u -v github.com/esimov/dithergo`

### Running

Type `go run cmd/main.go --help` to check all the supported commands. The library supports the following commands:

```
Usage of commands:
  -e string
    	Generates & exports the color and greyscale mode halftone images. 
	Options: 'all', 'color', 'mono' (default "all")
  -em float
    	Error multiplier (default 1.18)
  -o string
    	Output folder
  -t	Option to export the tresholded image (default true)

```
You can run all of the supported dithering algorithms at once, or you can run a specific one from the `cmd` directory.  

### Results:
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

## Author

* Endre Simo ([@simo_endre](https://twitter.com/simo_endre))

## License

Copyright © 2018 Endre Simo

This software is distributed under the MIT license found in the LICENSE file.
