package aof

import (
	"bufio"
	"github.com/blkcor/mini-redis/resp"
	"io"
	"os"
	"sync"
	"time"
)

// Aof represents an Append-Only File (AOF) for Redis-like databases
type Aof struct {
	File   *os.File
	Reader *bufio.Reader
	Mutex  sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		File:   f,
		Reader: bufio.NewReader(f),
	}
	// 开启一个go routine将aof同步到磁盘
	go func() {
		for {
			aof.Mutex.Lock()
			aof.File.Sync()
			aof.Mutex.Unlock()
			time.Sleep(time.Second * 1)
		}
	}()
	return aof, nil
}

// Close closes the AOF file
func (aof *Aof) Close() error {
	aof.Mutex.Lock()
	defer aof.Mutex.Unlock()
	return aof.File.Close()
}

// Write writes a RESP value to the AOF file
func (aof *Aof) Write(value resp.Value) error {
	aof.Mutex.Lock()
	defer aof.Mutex.Unlock()
	_, err := aof.File.Write(value.Marshal())
	if err != nil {
		return err
	}
	return nil
}

// Read reads RESP values from the AOF file
func (aof *Aof) Read(callback func(value resp.Value)) error {
	aof.Mutex.Lock()
	defer aof.Mutex.Unlock()

	resp := resp.NewResp(aof.File)
	for {
		value, err := resp.Read()
		if err == nil {
			callback(value)
		}
		if err == io.EOF {
			break
		}
		return err
	}
	return nil
}
