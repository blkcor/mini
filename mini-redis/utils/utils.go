package utils

import (
	"fmt"
	"github.com/blkcor/mini-redis/handler"
	"github.com/blkcor/mini-redis/resp"
	"strings"
)

// HandleCommand processes RESP commands
func HandleCommand(value resp.Value) (response resp.Value, err error) {
	command := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]

	h, ok := handler.Handler[command]
	if !ok {
		fmt.Println("Invalid command: ", command)
		return resp.Value{}, fmt.Errorf("invalid command: %s", command)
	}
	return h(args), nil
}
