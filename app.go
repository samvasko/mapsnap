package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	tiles := CreateMatrix(point{34550, 20755}, point{34555, 20760})
	fmt.Printf("Dimension x: %d, y: %d \n", tiles.width(), tiles.height())
	fmt.Printf("Number of tiles to download: %d\n", tiles.size())

	temp := initTemp()
	fmt.Println("\n Saving into temp directory: %s", temp)

	for dx := tiles.TL.x; dx < tiles.TR.x; dx++ {
		for dy := tiles.BL.y; dy < tiles.TL.y; dy++ {
			fmt.Printf("Getting tile: x: %d, y: %d ", dx, dy)
			get(dx, dy)
		}
	}

	Join(tiles)
}

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

func fetch(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}

	return resp.Body
}

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
