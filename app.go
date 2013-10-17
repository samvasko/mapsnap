package main

import (
	// "errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "mapsnap"
	app.Usage = "Download map tiles and join them"
	app.Flags = []cli.Flag{
		cli.StringFlag{"host", "mapbox", "Host for maps. Supported: mapbox,google,openstreetmap"},
		cli.IntFlag{"z", 16, "Z coordinate/Zoom factor"},
	}

	app.Action = func(c *cli.Context) {
		// Well defined rectangle
		if len(c.Args()) < 2 {
			bailout(c)
		}

		// Transform points into coordinates
		// Add better error handling
		tiles := CreateMatrix(
			digestPoint(c.Args()[0]),
			digestPoint(c.Args()[1]),
		)

		// Info time
		fmt.Printf("Dimension x: %d, y: %d \n", tiles.width(), tiles.height())
		fmt.Printf("Number of tiles to download: %d\n", tiles.size())

		// Move to temp and remember home
		script_dir, _ := os.Getwd()
		temp := initTemp()

		for dx := tiles.TL.x; dx < tiles.TR.x; dx++ {
			for dy := tiles.BL.y; dy < tiles.TL.y; dy++ {
				fmt.Printf("Getting tile: x: %d, y: %d ", dx, dy)
				get(dx, dy)
			}
		}

		wand := Join(tiles)
		err := os.Chdir(script_dir)
		handle(err)
		os.Remove(temp)

		wand.WriteImage(c.Args()[2])

	}
	app.Run(os.Args)

}

func bailout(c *cli.Context) {
	fmt.Println("Incorrect Usage.\n")
	cli.ShowAppHelp(c)
}

func digestPoint(strpoint string) point {
	splitstrpoint := strings.Split(strpoint, ",")
	var coords []coord
	for _, c := range splitstrpoint {
		nc, err := strconv.ParseUint(c, 10, 64)
		if err != nil {
			panic("Failied to parse point coordinates")
		}
		coords = append(coords, coord(nc))
	}

	return point{coords[0], coords[1]}
}

/**
 * Control downloading
 */
func get(x, y coord) {
	// Sample url: http://c.tile.openstreetmap.org//16/34551/20759.png
	// http://a.tile.openstreetmap.org//z/xxxx/yyyy.png

	url := fmt.Sprintf("http://a.tile.openstreetmap.org//16/%d/%d.png", x, y)
	filename := fmt.Sprintf("%d_%d.png", x, y)
	fmt.Print(url)

	file := fetch(url)
	save(file, filename)

	fmt.Print(" Done\n")
}

/**
 * Download and return reader
 */
func fetch(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		handle(err)
	}
	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	return resp.Body
}

/**
 * Save files
 */
func save(file io.ReadCloser, filename string) {
	out, err := os.Create(filename)
	handle(err)

	n, err := io.Copy(out, file)
	handle(err)
	_ = n
}

/**
 * Create temp dir and go there
 */
func initTemp() string {
	temp, err := ioutil.TempDir("", "mapTiles")
	handle(err)
	err = os.Chdir(temp)
	handle(err)

	return temp
}
