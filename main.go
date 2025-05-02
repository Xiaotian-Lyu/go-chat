package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. turn on the server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server is listening on port 8080")

	// 启动消息广播管理器
	go Manager()

	for {
		//2. wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		//3. handle the connection
		go handleConnection(conn)
	}
}

// func handleConnection(conn net.Conn) {
// 	addr := conn.RemoteAddr().String()
// 	fmt.Println("Client connected:", addr)
// 	defer func() {
// 		fmt.Println("Client disconnected:", addr)
// 		conn.Close()
// 	}()

// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			return
// 		}
// 		fmt.Printf("Received information from %s client: %s ", addr, string(buf[:n]))
// 	}
// }
