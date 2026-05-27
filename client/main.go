package main

import "flag"

func main() {
	var clientId string
	var serverAddr string

	flag.StringVar(&clientId, "client_id", "", "client id")
	flag.StringVar(&serverAddr, "server_addr", "", "server address")
	flag.Parse()

	client := NewClient(clientId, serverAddr)

	client.Run()
}
