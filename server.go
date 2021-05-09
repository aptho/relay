package main

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
	ips []*url.URL
}

type requestHandler struct {
	server Server
}

func Setup(ips []string) Server {
	urls := make([]*url.URL, 0)
	for _, ip := range ips {
		url, err := url.Parse(ip)
		if err != nil {
			fmt.Println(err)
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

	for _, ip := range handler.server.ips {
		go sendRequest(ip, req, body)
	}
}

func sendRequest(ip *url.URL, r *http.Request, body bytes.Buffer) {
	req := r.WithContext(context.TODO())
	req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
	req.URL = ip
	req.URL.Path = r.URL.String()
	req.RequestURI = ""

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
