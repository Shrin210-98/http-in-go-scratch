package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// buf := make([]byte, 1024)
	// conn.Read(buf)
	// req := string(buf)
	// fmt.Println("URL", strings.Join(strings.Split(req, " "), " II "))

	reader := bufio.NewReader(conn)
	tcpPayload, _ := reader.ReadString('\n')
	fmt.Println("Data: ", tcpPayload)

	// methods := []string{"GET","POST","PUT","DELETE"}

	for index, item := range strings.Split(tcpPayload, " ") {
		if index == 1 {
			if item == "/" {
				conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			} else {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			}
		}
	}

}
