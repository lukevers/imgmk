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

	// Get the physical image ready with colors/size
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, *width, *height))

	// If we passed the transparent flag then we want the
	// background to be transparent.
	if *transparent {
		bg = image.Transparent
	}

	// Draw the empty image
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	// Create new freetype context to get ready for
	// adding text.
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(font)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(freetype.NoHinting)

	// If we're piping in text then we might have multiple lines,
	// and we'll have to draw each line separately.
	if fi.Mode()&os.ModeNamedPipe != 0 {
		reader := bufio.NewReader(os.Stdin)
		var count float64 = 1
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				// When we get here it means we're out of lines,
				// so we just break out of here instead of reporting
				// any errors.
				break
			}

			drawText(c, line, *size, count)
			count++
		}
	} else {
		// If we're down here we're just doing what I assume to be one line,
		// which was passed in either not at all (so using default), or by
		// the *text flag.
		drawText(c, *text, *size, 1)
	}

	// Create image
	image, err := os.Create(*name)
	if err != nil {
		fmt.Printf("Error creating file: %s", err)
		os.Exit(1)
	}

	// Defer closing of image until we're done, but make sure we actually
	// close it.
	defer image.Close()

	// Write the image and encode it to png.
	b := bufio.NewWriter(image)
	err = png.Encode(b, rgba)
	if err != nil {
		fmt.Printf("Error encoding image: %s", err)
		os.Exit(1)
	}

	// This actually writes all the data to the image.
	err = b.Flush()
	if err != nil {
		fmt.Printf("Error flushing bufio writer: %s", err)
		os.Exit(1)
	}
}

func drawText(c *freetype.Context, text string, size, line float64) {
	// We need an offset because we need to know where exactly on the
	// image to place the text. The `line` is how much of an offset
	// that we need to provide (which line the text is going on).
	offsetY := 10 + int(c.PointToFix32(size*line)>>8)

	_, err := c.DrawString(text, freetype.Pt(10, offsetY))
	if err != nil {
		fmt.Println("Could not draw text: %s", err)
		os.Exit(1)
	}
}
