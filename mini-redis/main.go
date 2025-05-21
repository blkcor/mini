package main

import (
	"fmt"
	"github.com/blkcor/mini-redis/handler"
	"github.com/blkcor/mini-redis/resp"
	"net"
	"strings"
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
		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]
		handle, ok := handler.Handler[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			respWriter.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}
		respWriter.Write(handle(args))
	}
}
