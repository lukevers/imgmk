package main

import (
	"github.com/lukevers/freetype-go/freetype"
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
	"flag"
	"image"
	"image/draw"
	"image/png"
)

func main() {
	// Parse flags
	flag.Parse()

	// Load the font file
	file, err := ioutil.ReadFile(*path)
	if err != nil {
		fmt.Printf("Error loading font from file: %s", err)
		os.Exit(1)
	}

	// Now parse the ttf font
	font, err := freetype.ParseFont(file)
	if err != nil {
		fmt.Printf("Error parsing font: %s", err)
		os.Exit(1)
	}

	// Draw image
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, *width, *height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(font)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(freetype.NoHinting)

	// Draw text
	_, err = c.DrawString(*text, freetype.Pt(10, 10+int(c.PointToFix32(*size)>>8)))

	// Create image
	image, err := os.Create(*name)
	if err != nil {
		fmt.Printf("Error creating file: %s", err)
		os.Exit(1)
	}

	defer image.Close()

	b := bufio.NewWriter(image)
	err = png.Encode(b, rgba)
	if err != nil {
		fmt.Printf("Error encoding image: %s", err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		fmt.Printf("Error flushing bufio writer: %s", err)
		os.Exit(1)
	}


}
