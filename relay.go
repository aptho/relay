package relay

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Server struct {
	backends []*url.URL
}

type requestHandler struct {
	server Server
}

func Setup(backends []string) Server {
	// Create slice to hold valid urls
	urls := make([]*url.URL, 0)

	for _, backend := range backends {
		// Attempt to parse the backend url
		url, err := url.Parse(backend)
		if err != nil {
			fmt.Printf("Invalid backend url: %s\n", backend)
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}

		urls = append(urls, url)
	}

	return Server{urls}
}

func (s Server) Start(port string) {
	handler := requestHandler{s}
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (handler requestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Pull the body from the request because it drains after the first read
	var body bytes.Buffer
	body.ReadFrom(req.Body)

	for _, backend := range handler.server.backends {
		// For each backend, pass the request on
		go sendRequest(backend, req, body)
	}
}

func sendRequest(backend *url.URL, r *http.Request, body bytes.Buffer) {
	// Clone the request with an empty context
	req := r.WithContext(context.TODO())

	// Add back the body, and set the new request url as the current backend
	req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
	req.URL = backend
	req.URL.Path = r.URL.String()
	req.RequestURI = ""

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
