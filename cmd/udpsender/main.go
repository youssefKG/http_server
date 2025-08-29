package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		fmt.Print("error creating udp resolver")

		return
	}

	fmt.Println("Resolved udp add: ", raddr)

	conn, err := net.DialUDP("udp", nil, raddr)

	if err != nil {
		fmt.Println("cat create dial udp", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("UDP sender started. Type messages and press Enter.")
	fmt.Println("Press Ctrl+C to quit.")

	for {
		buff := make([]byte, 16)
		n, err := reader.Read(buff)
		fmt.Print(">")
		if err != nil {
			fmt.Println("Error reading from input")
			continue
		}
		conn.Write(buff[:n])
	}

}
