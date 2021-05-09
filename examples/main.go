package main

import "github.com/aptho/relay"

func main() {
	server := relay.Setup([]string{"http://localhost:8081", "http://localhost:8082"})
	server.Start("8080")
}
