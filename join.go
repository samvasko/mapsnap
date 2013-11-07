package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"os"
)

func Join(tiles matrix) image.Image {
	cols := []string{}

	// Itinerate over Xes
	for dx := tiles.TL.x; dx < tiles.TR.x; dx++ {
		col := []string{}
		// All y for this x
		for dy := tiles.BL.y; dy < tiles.TL.y; dy++ {
			col = append(col, fmt.Sprintf("tile_%d:%d.png", dx, dy))
		}
		img := joinByName(col, true)
		filename := fmt.Sprintf("column_%d.png", dx)
		fmt.Print("Created: ", filename, "                   \r")
		Save(img, filename)
		cols = append(cols, filename)
	}

	return joinByName(cols, false)
}

func Save(img image.Image, filename string) {
	outfile, err := os.Create(filename)
	handle(err)
	png.Encode(outfile, img)
	outfile.Close()
}

func joinByName(filenames []string, vertical bool) image.Image {

	var output *image.RGBA
	var tile_size struct{ dx, dy int }
	var offset int

	for i, fname := range filenames {
		tile := open(fname)
		if i == 0 {
			offset = 0
			tile_size.dx = tile.Bounds().Dx()
			tile_size.dy = tile.Bounds().Dy()

			if vertical {
				output = createHost(tile_size.dx, tile_size.dy*len(filenames))
			} else {
				output = createHost(tile_size.dx*len(filenames), tile_size.dy)
			}
		}

		for dx := 0; dx < tile_size.dx; dx++ {
			for dy := 0; dy < tile_size.dy; dy++ {
				if vertical {
					output.Set(dx, dy+(offset*tile_size.dy), tile.At(dx, dy))
				} else {
					output.Set(dx+(offset*tile_size.dx), dy, tile.At(dx, dy))
				}
			}
		}

		offset++
	}

	return output
}

func createHost(xb, yb int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, xb, yb))
}

func open(filename string) image.Image {
	file, err := os.Open(filename)
	handle(err)
	img, _, err := image.Decode(file)
	handle(err)

	file.Close()

	img = resize.Resize(0, 150, img, resize.NearestNeighbor)
	return img
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
