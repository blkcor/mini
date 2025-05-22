package main

import (
	"fmt"
	"github.com/blkcor/mini-redis/aof"
	"github.com/blkcor/mini-redis/resp"
	"github.com/blkcor/mini-redis/utils"
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
	aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()
	// 监听客户端连接
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	// 从.aof文件中读取数据
	err = aof.Read(func(value resp.Value) {
		_, err := utils.HandleCommand(value)
		if err != nil {
			return
		}
	})
	if err != nil {
		fmt.Println(err)
		return
	}
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
		response, err := utils.HandleCommand(value)
		if err != nil {
			respWriter.Write(resp.Value{Typ: "string", Str: ""})
			continue
		}
		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}
		respWriter.Write(response)
	}
}
