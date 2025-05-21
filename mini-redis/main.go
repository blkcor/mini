package main

import (
	"fmt"
	"github.com/blkcor/mini-redis/resp"
	"net"
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
		respReader := resp.NewResp(conn)
		respWriter := resp.NewWriter(conn)

		value, err := respReader.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		respWriter.Write(resp.Value{
			Typ: "string",
			Str: "OK",
		})
	}
}
