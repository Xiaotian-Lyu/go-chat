# Go Chat Project

A simple chat system written in Go, including both server and client.

## 🚀 Features

- Multi-user chat room (`who` command to list online users)
- Rename username (`rename|new_name`)
- Private messaging (`to|username|message`)
- Auto kick on inactivity (disconnect after 60 seconds of inactivity)
- Simple TCP client (`client.go`)

## 🛠️ How to Run

1. Run the server

    In the project root directory, run:
    
    go run server.go manager.go user.go

2. Run the client

    In another terminal window, run:
    
    go run client.go
    
    You can launch multiple clients in separate terminal windows for testing.

## 💬 Command Examples

List online users: who

Change username: rename|new_name

Send a private message: to|username|message

## ⚠️ Notes

Do not run `go run .` directly, as it will compile both `client.go` and `server.go` and cause a `main()` conflict error.

Always start the server and client in separate terminal windows.

If using VSCode, you can quickly open the terminal with Ctrl + ~.

## ✨ Future Improvements

- Web UI client
- Database support for chat history
- User authentication and login system
