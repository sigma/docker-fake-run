package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/njern/httpstream"
)

func getStream(url string) (chan string, chan error) {
	// fmt.Println(url)
	stream := make(chan string)
	client := httpstream.NewClient(func(line []byte) {
		stream <- string(line)
	})
	done := make(chan error)
	err := client.Connect(url, done)
	if err != nil {
		panic(err)
	}

	return stream, done
}

func getStreamUrl(base string, output string) string {
	values := url.Values{}
	values.Add(output, "1")
	values.Add("follow", "1")
	return base + "?" + values.Encode()
}

func readStream(url string, output io.Writer) <-chan interface{} {
	stream, done := getStream(url)
	_done := make(chan interface{})

	go func() {
		for {
			select {
			case event := <-stream:
				fmt.Fprintf(output, event)
			case <-done:
				_done <- true
				return
			}
		}
	}()
	return _done
}

type Status struct {
	StatusCode int
}

func main() {
	base := os.Args[1]
	container_logs := base + "/logs"
	container_wait := base + "/wait"

	stdout_url := getStreamUrl(container_logs, "stdout")
	stdout_done := readStream(stdout_url, os.Stdout)

	stderr_url := getStreamUrl(container_logs, "stderr")
	stderr_done := readStream(stderr_url, os.Stderr)

	<-stdout_done
	<-stderr_done

	resp, _ := http.Post(container_wait, "", nil)
	status := new(Status)
	json.NewDecoder(resp.Body).Decode(&status)

	os.Exit(status.StatusCode)
}
