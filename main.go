package main

func main() {
	server := NewApiServer(":2222")
	server.Start()
}
