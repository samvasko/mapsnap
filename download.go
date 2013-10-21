package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
	"time"
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
		var dfile io.ReadCloser

		// Download with increasing time
		for i := 0; i < 10; i++ {
			file, err := fetch(url)
			if err != nil {
				if i < 8 {
					fmt.Println(err)
					takeNap(i)
				} else {
					bailout(err)
				}
			} else {
				dfile = file
				break
			}
		}

		for i := 0; i < 10; i++ {
			err := save(dfile, j.toFilename())
			if err != nil {
				if i < 8 {
					fmt.Println(err)
					takeNap(i)
				} else {
					bailout(err)
				}
			} else {
				fmt.Print(j.toFilename(), "\r")
				break
			}
		}

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
func fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return resp.Body, nil
}

/**
 * Save files
 */
func save(file io.ReadCloser, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, file)

	if err != nil {
		return (err)
	}

	file.Close()

	return nil
}

func takeNap(t int) {
	fmt.Println("Problem ocurred, taking a nap")
	time.Sleep(time.Second * 10 * time.Duration(t))
}
