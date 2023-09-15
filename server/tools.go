package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func getName(conn net.Conn) (string, error) {
	name, _ := bufio.NewReader(conn).ReadString('\n')
    name = strings.TrimSpace(name)
	if len(name)==0 ||len(name)>20{
		return "", fmt.Errorf("the length of the name must be at least one character and less than 20 characters")
	}
	return name, nil
}

func printPenguin(conn net.Conn){
	for _, v := range pinguin {
		fmt.Fprintln(conn, v)
	}
}

func checkMsg(msg string) bool{
	for _, v := range msg{
		if v >= '!' && v <= '~' {
			return true
		}
	}
	return false
}

func checkLen(msg string) bool{
	if len(msg) > 0 && len(msg) <= 200{
		return true
	}
	return false
}