package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"netcat/messenger"
)

// Default settings
const (
	Port           = ":8989"
	MaxConnections = 10
)

func main() {
	port, err := GetPort()
	if err != nil {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}
	server := &messenger.Server{}
	err = server.Constructor(port, MaxConnections)
	if err != nil {
		log.Printf("ERROR -> main: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Listening on the port %v\n", port)
	closeSignal := make(chan os.Signal, 1)
	go CloseProgramOnSignal(server, closeSignal)
	for {
		conn, err := server.Server.Accept()
		if err != nil {
			break
		}
		go server.ConnectMessenger(conn)
	}
	<-closeSignal
}

// GetPort - Gets Port from user args
func GetPort() (string, error) {
	args := os.Args
	if len(args) < 2 {
		return Port, nil
	} else if len(args) > 2 {
		return "", errors.New("User inputs more than 2 args")
	}
	return ":" + os.Args[1], nil
}

// CloseProgramOnSignal - Closing program on close signal input
func CloseProgramOnSignal(server *messenger.Server, closeSignal chan os.Signal) {
	signal.Notify(closeSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// Wait signal
	<-closeSignal
	server.CloseServer()
	os.Exit(0)
}
