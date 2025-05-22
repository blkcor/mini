package handler

import (
	"github.com/blkcor/mini-redis/resp"
	"sync"
)

// Handler mapping of command to handler function
var Handler = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hSet,
	"HGET": hGet,
}

var Sets = map[string]string{}
var SetsMu = sync.RWMutex{}

var HSets = map[string]map[string]string{}
var HSetsMu = sync.RWMutex{}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}
	return resp.Value{Typ: "string", Str: args[0].Bulk}
}

// set command handler
func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].Bulk
	value := args[1].Bulk
	SetsMu.Lock()
	Sets[key] = value
	SetsMu.Unlock()
	return resp.Value{Typ: "string", Str: "OK"}
}

// get command handler
func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}
	SetsMu.RLock()
	defer SetsMu.RUnlock()
	key := args[0].Bulk
	value, ok := Sets[key]
	if !ok {
		return resp.Value{Typ: "null"}
	}
	return resp.Value{Typ: "bulk", Bulk: value}
}

// hSet command handler
func hSet(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}
	key := args[0].Bulk
	field := args[1].Bulk
	value := args[2].Bulk
	HSetsMu.Lock()
	if _, ok := HSets[key]; !ok {
		HSets[key] = make(map[string]string)
	}
	HSets[key][field] = value
	HSetsMu.Unlock()
	return resp.Value{Typ: "string", Str: "OK"}
}

// hGet command handler
func hGet(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}
	HSetsMu.RLock()
	defer HSetsMu.RUnlock()
	key := args[0].Bulk
	field := args[1].Bulk
	fields, ok := HSets[key]
	if !ok {
		return resp.Value{Typ: "null"}
	}
	value, ok := fields[field]
	if !ok {
		return resp.Value{Typ: "null"}
	}
	return resp.Value{Typ: "bulk", Bulk: value}
}
