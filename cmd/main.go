package main

import (
	"github.com/esimov/dithergo"
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"FloydSteinberg",
		dither.Settings{
			[][]float32{
				[]float32{ 0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0 },
				[]float32{ 3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0 },
				[]float32{ 1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0 },
			},
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
		},
	},
}

func main() {
	dither.Process(ditherers)
}