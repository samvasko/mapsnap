package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"time"
)

var hosts = map[string]string{
	"mapbox":        "",
	"openstreetmap": "http://a.tile.openstreetmap.org//{{.Z}}/{{.X}}/{{.Y}}.png",
	"googlemaps":    "",
}

func main() {
	app := cli.NewApp()
	app.Name = "mapsnap"
	app.Usage = "Download map tiles and join them"
	app.Flags = []cli.Flag{
		cli.StringFlag{"host", "mapbox", "Host for maps. Supported: mapbox,google,openstreetmap"},
		cli.IntFlag{"z", 0, "Z coordinate/Zoom factor"},
	}

	app.Action = func(c *cli.Context) {

		points, filename := parseArgs(c)
		z, host := parseFlags(c)

		// Transform points into coordinates and create
		// Value matrix
		tiles := CreateMatrix(points)

		// Info time
		fmt.Printf("Dimension x: %d, y: %d \n", tiles.width(), tiles.height())
		fmt.Printf("Number of tiles to download: %d\n", tiles.size())

		// Move to temp and remember home
		script_dir, _ := os.Getwd()
		temp := initTemp()

		// Create job channel
		jobs := make(chan Downloadjob, int(tiles.size()))
		done := make(chan bool, int(tiles.size()))

		// Create workers
		for i := 0; i < 4; i++ {
			go Downloader(i, jobs, done)
		}

		// Send jobs
		for dx := tiles.TL.x; dx < tiles.TR.x; dx++ {
			for dy := tiles.BL.y; dy < tiles.TL.y; dy++ {
				jobs <- Downloadjob{point{dx, dy}, coord(z), host, 0}
			}
		}

		close(jobs)

		// count finished
		start := time.Now()

		for i := 0; i < int(tiles.size()); i++ {
			<-done
			fmt.Printf("Finished: %d ~~ Speed: %.3f tiles/second \r", i, float64(i)/time.Since(start).Seconds())
		}

		img := Join(tiles)
		os.Chdir(script_dir)
		Save(img, filename)
		os.Remove(temp)

	}

	app.Run(os.Args)

}

func bailout(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func parseArgs(c *cli.Context) ([2]string, string) {
	// Well defined rectangle
	if len(c.Args()) < 2 {
		bailout(errors.New("Incorrect point definition: Define points like this '10,10 20,34'"))
		cli.ShowAppHelp(c)
	}

	// Output file
	var outfile string
	if len(c.Args()) < 3 {
		outfile := fmt.Sprintf("__%s_%s.png", c.Args()[0], c.Args()[1])
		fmt.Printf("Output filename not specified using %s", outfile)
	} else {
		outfile = c.Args()[2]
	}

	return [2]string{c.Args()[0], c.Args()[1]}, outfile
}

func parseFlags(c *cli.Context) (int, string) {
	var z int
	if c.Int("z") != 0 {
		z = c.Int("z")
	} else {
		bailout(errors.New("Missing z"))
		cli.ShowAppHelp(c)
		z = 0
	}

	host := "openstreetmap"
	return z, host
}

/**
 * Create temp dir and go there
 */
func initTemp() string {
	temp, err := ioutil.TempDir("", "mapTiles")
	if err != nil {
		bailout(err)
	}

	os.Chdir(temp)
	return temp
}
