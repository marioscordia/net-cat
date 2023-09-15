package main

import (
	"fmt"
	"netcat/server"
	"os"
)
func main(){
	args := os.Args

	if len(args) == 1 {
		server.StartServer("8989")
    }else if len(args) == 2 {
		server.StartServer(args[1])
	}else{
		fmt.Println("[USAGE]: go run . $port <- (optional)")
		os.Exit(1)
	}

}