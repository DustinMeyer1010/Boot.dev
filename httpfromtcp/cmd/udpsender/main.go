package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	raddr, err := net.ResolveUDPAddr("udp", ":42069")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)

	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")

		readString, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte(readString))
	}

}
