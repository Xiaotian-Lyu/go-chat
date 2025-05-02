package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string      //客户端的 IP:Port，默认当作唯一ID
	Conn net.Conn    //可以用它 Conn.Write() 或 Conn.Read() 直接操作客户端
	Chan chan string //服务器给这个用户发消息的专用通道
}

func (u *User) ListenMessage() {
	for msg := range u.Chan {
		u.Conn.Write([]byte(msg))
	}
}

func (u *User) DoMessage(msg string) {
	cleanMsg := strings.TrimSpace(msg)
	switch {
	case cleanMsg == "who":
		userLock.Lock()
		for _, user := range onlineUsers {
			u.Chan <- "[" + user.Addr + "] " + user.Name + "\n"
		}
		userLock.Unlock()
	case strings.HasPrefix(cleanMsg, "rename|"):
		parts := strings.SplitN(cleanMsg, "|", 2)
		if len(parts) < 2 {
			u.Chan <- "格式错误: rename|新名字\n"
			return
		}
		newName := parts[1]
		userLock.Lock()
		if _, ok := onlineUsers[newName]; ok {
			u.Chan <- "名字已被占用\n"
		} else {
			delete(onlineUsers, u.Name)
			u.Name = newName
			onlineUsers[u.Name] = u
			u.Chan <- "用户名修改成功: " + u.Name + "\n"
		}
		userLock.Unlock()
	default:
		message <- "[" + u.Name + "]: " + cleanMsg
	}
}
