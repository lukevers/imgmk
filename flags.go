package main

import (
	"flag"
)

var (
	path        = flag.String("font-path", "fonts/itchy.ttf", "Path to font file")
	name        = flag.String("name", "out.png", "Image name output.")
	text        = flag.String("text", "Test", "The text to be displayed")
	size        = flag.Float64("size", 16, "The font size to use")
	width       = flag.Int("width", 640, "The width of the image")
	height      = flag.Int("height", 480, "The height of the image")
	dpi         = flag.Float64("dpi", 300, "The DPI of the image.")
	transparent = flag.Bool("transparent", false, "Transparent image background")
)
