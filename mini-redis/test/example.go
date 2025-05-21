package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Test() {
	input := "$5\r\nhello\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	b, _ := reader.ReadByte()
	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}
	size, _ := reader.ReadByte()
	strSize, _ := strconv.ParseInt(string(size), 10, 64)
	// 消费 \r\n
	reader.ReadByte()
	reader.ReadByte()

	content := make([]byte, strSize)
	reader.Read(content)

	fmt.Println("Read content:", string(content))
}
