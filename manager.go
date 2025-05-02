package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

//多 goroutine 传消息 → channel
//多 goroutine 改同一个资源 → lock

var (
	onlineUsers = make(map[string]*User) //保存在线用户的 map，key 是地址（IP:port），value 是 *User 对象。
	userLock    sync.Mutex               //保护 onlineUsers 的并发安全，防止多个 goroutine 同时读写崩溃
	message     = make(chan string)
)

func Manager() {
	for {
		msg := <-message //从 message 这个 channel 中取出一条消息，
		userLock.Lock()
		for _, user := range onlineUsers { //key 是地址，value 是 *User
			user.Chan <- msg + "\n"
		}
		userLock.Unlock()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String() //获取客户端的地址，格式是 IP:port
	user := &User{
		Name: addr,
		Addr: addr,
		Conn: conn,
		Chan: make(chan string),
	}

	userLock.Lock()
	onlineUsers[user.Name] = user //用户作为index
	userLock.Unlock()

	// 通知所有人上线消息
	message <- fmt.Sprintf("[%s] 上线了", user.Name)

	// 启动监听当前用户发送消息的 goroutine
	go user.ListenMessage()

	isActive := make(chan bool)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			msg := string(buf[:n])
			user.DoMessage(msg)
			isActive <- true // 每收到一条消息，告诉主 goroutine 用户活跃
		}

		// 用户断开时清理
		userLock.Lock()
		delete(onlineUsers, addr)
		userLock.Unlock()
		message <- fmt.Sprintf("[%s] 下线了", user.Name)
		close(user.Chan)
		close(isActive)   // 关闭 isActive channel
		isActive <- false // 用户断开连接，发送一个信号表示不活跃
	}()

	for {
		select {
		case active, ok := <-isActive:
			if !ok || !active {
				return
			}
		case <-time.After(60 * time.Second):
			user.Chan <- "超时未操作，您已被踢出\n"
			conn.Close() // 主动断开连接
			return
		}
	}

}
