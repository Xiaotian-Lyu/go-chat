package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    // 连接服务器
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("连接失败:", err)
        return
    }
    defer conn.Close()

    // 启动一个 goroutine 专门接收服务器消息
    go func() {
        buf := make([]byte, 1024)
        for {
            n, err := conn.Read(buf)
            if err != nil {
                fmt.Println("服务器关闭连接")
                os.Exit(0)
            }
            fmt.Print(string(buf[:n]))
        }
    }()

    // 主 goroutine 读控制台输入并发给服务器
    input := bufio.NewReader(os.Stdin)
    for {
        line, err := input.ReadString('\n')
        if err != nil {
            fmt.Println("输入出错:", err)
            break
        }
        conn.Write([]byte(line))
    }
}
