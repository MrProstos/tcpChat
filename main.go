package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
)

func main() {
	l, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return
	}

	defer l.Close()

	connMap := &sync.Map{}

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}

		conn.Write([]byte("Введите свое имя\n"))
		name, _ := bufio.NewReader(conn).ReadString('\n')

		id := uuid.New().String()
		connMap.Store(id, conn)

		go HadleConnect(conn, connMap, id, name[:len(name)-1])
	}
}

func HadleConnect(conn net.Conn, connMap *sync.Map, id string, name string) {
	defer func() {
		conn.Close()
		connMap.Delete(id)
	}()

	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}

		connMap.Range(func(key, value any) bool {
			if c, ok := value.(net.Conn); ok {
				if conn != value {
					if _, err := c.Write([]byte(fmt.Sprintf("%v:%v", name, userInput))); err != nil {
						log.Println(err)
					}
				}
			}

			return true
		})
	}
}
