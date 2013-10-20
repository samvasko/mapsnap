package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

type Downloadjob struct {
	target point
	z      coord
	host   string
}

func Downloader(id int, jobs <-chan Downloadjob) {
	fmt.Println("Initialized worker", id)
	for j := range jobs {
		url := j.toUrl()
		fmt.Println("Downloading url: ", url)
		file := fetch(url)
		save(file, j.toFilename())
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
func fetch(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		bailout(err)
	}
	if resp.StatusCode != 200 {
		err := fmt.Sprintf("boop, encoutered error while checking out: %s \n status: %s ", url, resp.Status)
		bailout(errors.New(err))
	}

	return resp.Body
}

/**
 * Save files
 */
func save(file io.ReadCloser, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		bailout(err)
	}
	n, _ := io.Copy(out, file)
	_ = n
}
