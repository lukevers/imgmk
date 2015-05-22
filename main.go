package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/lukevers/freetype-go/freetype"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
)

func main() {
	// Parse flags
	flag.Parse()

	// Check to see if we're passing text in via a unix pipe
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error running stat on stdin: %s", err)
		os.Exit(1)
	}

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
	c.SetDPI(*dpi)
	c.SetFont(font)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(freetype.NoHinting)

	// Draw text
	if fi.Mode()&os.ModeNamedPipe != 0 {
		reader := bufio.NewReader(os.Stdin)
		var count float64 = 1
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			drawText(c, line, *size, count)
			count++
		}
	} else {
		// Draw text
		drawText(c, *text, *size, 1)
	}

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

func drawText(c *freetype.Context, text string, size, line float64) {
	offsetY := 10 + int(c.PointToFix32(size*line)>>8)

	_, err := c.DrawString(text, freetype.Pt(10, offsetY))
	if err != nil {
		fmt.Println("Could not draw text: %s", err)
		os.Exit(1)
	}
}
