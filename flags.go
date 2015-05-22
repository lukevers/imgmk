package main

import (
	"flag"
)

var path = flag.String("font-path", "fonts/itchy.ttf", "Path to font file")
var name = flag.String("name", "out.png", "Image name output.")
var text = flag.String("text", "Test", "The text to be displayed")
var size = flag.Float64("size", 16, "The font size to use")
var width = flag.Int("width", 640, "The width of the image")
var height = flag.Int("height", 480, "The height of the image")
var dpi = flag.Float64("dpi", 300, "The DPI of the image.")
