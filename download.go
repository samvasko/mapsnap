package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	// "runtime"
	"text/template"
	"time"
)

type Downloadjob struct {
	target  point
	z       coord
	host    string
	penalty int
}

var (
	finished  uint64  = 0
	totaltime float64 = 0
)

func Downloader(id int, jobs chan Downloadjob, done chan bool) {
	fmt.Println("Initialized worker", id)
	for j := range jobs {
		err, _ := fetch(j.toUrl(), j.toFilename())
		if err != nil {
			fmt.Println("Error downloading", j.target, "With: ", err)
			if j.penalty > 20 {
				fmt.Print("Can't take this anymore")
				panic(err)
			}

			j.penalty++
			takeNap(j.penalty)
		}

		// Yay success
		done <- true
		// fmt.Printf("Current average: %.4f  -- Routines: %d \r", runtime.NumGoroutine())
	}
}

func (d *Downloadjob) toUrl() string {
	urltemplate := hosts[d.host]
	t := template.Must(template.New("url").Parse(urltemplate))

	// Parse template to file
	var out bytes.Buffer
	err := t.Execute(&out, struct{ X, Y, Z coord }{d.target.x, d.target.y, d.z})
	if err != nil {
		bailout(err)
	}
	// Return created url
	return out.String()
}

func (d *Downloadjob) toFilename() string {
	return fmt.Sprintf("tile_%d:%d.png", d.target.x, d.target.y)
}

/**
 * Download and return reader
 */
func fetch(url, filename string) (error, float64) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return err, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status), 0
	}
	out, err := os.Create(filename)
	if err != nil {
		return err, 0
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err, 0
	}

	return nil, time.Since(start).Seconds()
}

func takeNap(t int) {
	delay := time.Second * 5 * time.Duration(t)
	fmt.Println("Problem ocurred, taking a nap for", delay, " sec")
	time.Sleep(delay)
}
