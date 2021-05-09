package main

func main() {
	server := Setup([]string{"http://localhost:8081", "http://localhost:8082"})
	server.Start("8080")
}
