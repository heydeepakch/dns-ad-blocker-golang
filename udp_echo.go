package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. Define the address to listen on
	addr := net.UDPAddr{
		Port: 9000,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// 2. Start listening on UDP
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("UDP Echo server listening on 127.0.0.1:9000")

	// 3. Create a buffer to hold incoming data
	buffer := make([]byte, 1024)

	for {
		// 4. Read incoming UDP packet
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("Received from %s: %s\n", clientAddr, message)

		// 5. Send response back to sender
		_, err = conn.WriteToUDP(buffer[:n], clientAddr)
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}