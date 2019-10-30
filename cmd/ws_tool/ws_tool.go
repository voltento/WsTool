package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/voltento/WsTool/internal/utils"
	"github.com/voltento/WsTool/pkg/command"
	"github.com/voltento/WsTool/pkg/handlers"
	"github.com/voltento/WsTool/pkg/web_socket_client"
)

func main() {

	address, headers, commandsFilePath := utils.ParseArgs()

	ws := new(web_socket_client.WebSocket)
	if er := ws.Connect(string(address), headers); er != nil {
		println("Error on connection. Reason: " + er.Error())
		os.Exit(1)
	}

	if len(commandsFilePath) > 0 {
		mgr := handlers.CreateHandlerManager(ws)
		for cmd := range command.CreateReaderFromFile(commandsFilePath) {
			if er := mgr.Handle(cmd); er != nil {
				fmt.Printf("Error occured. Error: %v", er.Error())
				os.Exit(1)
			}
		}
	}

	go printMessageFromWs(ws)

	readFromConsoleAndSendToWs(ws)
}

func printMessageFromWs(ws *web_socket_client.WebSocket) {
	var err error
	var msg string
	for {
		msg, err = ws.ReadOneMessage()
		if err != nil {
			fmt.Printf("Error occurred during read from socket. Error: %s\n", err.Error())
			os.Exit(0)
		}
		print("< : ", msg, "\n")
	}
}

func readFromConsoleAndSendToWs(ws *web_socket_client.WebSocket) {
	for {
		reader := bufio.NewReader(os.Stdin)

		var messageToWs string
		var err error

		messageToWs, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Can't read from console. Error: %s\n", err.Error())
		}

		err = ws.SendMessage(messageToWs)
		if err != nil {
			panic("Error occurred during send message to ws. Error: " + err.Error())
		} else {
			fmt.Printf(" > : %s\n", messageToWs)
		}
	}
}