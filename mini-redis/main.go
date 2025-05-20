package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listen on port :6379")
	// 创建一个server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 监听客户端连接
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		// 读取客户端发送的数据
		_, err = conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}
		// 忽略请求 直接返回 +OK
		conn.Write([]byte("+OK\r\n"))
	}
}
