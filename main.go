package main

import (
	"fmt"
	"strings"

	"github.com/Surya7890/ws-go/server"
	"github.com/Surya7890/ws-go/ws"
)

func main() {
	server := server.NewServer(":7000")
	go func() {
		for msg := range server.Msg {
			str := string(msg.Data)
			arr := strings.Split(str, " ")
			fmt.Println(arr[0], arr[1], msg.Sender)
			switch arr[0] {
			case "dm":
				err := ws.SendPrivateMessage(msg.Sender, arr[1], arr[2:])
				if err != nil {
					ws.ActivePeers[msg.Sender].Write([]byte(fmt.Sprintf("%s is offline\n", arr[1])))
				}
			}
		}
	}()
	server.Start()
}
