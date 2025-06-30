package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnections(conn)
	}
}

func handleConnections(conn net.Conn) {

	defer conn.Close()

	// buf := make([]byte, 1024)
	// conn.Read(buf)
	// req := string(buf)
	// fmt.Println("URL", strings.Join(strings.Split(req, " "), " II "))

	reader := bufio.NewReader(conn)
	tcpPayload, _ := reader.ReadString('\n')
	tcpHost, _ := reader.ReadString('\n')
	tcpAccept, _ := reader.ReadString('\n')
	tcpUserAgent, _ := reader.ReadString('\n')
	tcpData := fmt.Sprintf("%s\r%s\r%s\r%s\r\n", tcpPayload, tcpHost, tcpAccept, tcpUserAgent)
	fmt.Println("TCP Request:\n", tcpData)

	// methods := []string{"GET","POST","PUT","DELETE"}

	for index, item := range strings.Split(tcpPayload, " ") {
		if index == 1 {
			if item == "" {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			} else {
				var httpPayload = ""
				var (
					httpStatus    = "HTTP/1.1 200 OK"
					request       = strings.Split(item, "/")
					contentType   = "text/plain"
					contentLength = ""
					response      = ""
				)
				switch request[1] {
				case "echo":
					contentLength = strconv.Itoa(len(request[2]))
					response = request[2]
					httpPayload = fmt.Sprintf(
						"%s\r\nContent-Type: %s\r\nContent-Length: %s\r\n\r\n%s",
						httpStatus,
						contentType,
						contentLength,
						response,
					)
					// else if path[:5] == "/echo" {
					// dynamic_path := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
					// res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(dynamic_path), dynamic_path)
				case "user-agent":
					response = strings.TrimSpace(strings.TrimPrefix(tcpUserAgent, "User-Agent: "))
					contentLength = strconv.Itoa(len(response))
					httpPayload = fmt.Sprintf(
						"%s\r\nContent-Type: %s\r\nContent-Length: %s\r\n\r\n%s",
						httpStatus,
						contentType,
						contentLength,
						response,
					)
				default:
					httpPayload = fmt.Sprintf("%s\r\n\r\n", httpStatus)
				}
				conn.Write([]byte(httpPayload))
			}
		}
	}
}
