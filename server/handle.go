package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)


var (
	Messages = make(chan Message)
	welcomeMessages = make(chan User)
	quitMessages = make(chan User)
)


func StartServer(port string){
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Printf("Server started, listening on the port %s\n", port)
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		if len(clients)==10{
			fmt.Fprint(conn, "Too many clients, disconnecting")
            conn.Close()
			continue
		}
		go handleConn(conn)
	}
}
// \n\033[1A\033[K

func broadcaster() {
	for{
        select {
		case inUser := <-welcomeMessages:
			for _, c := range clients {
				if inUser.conn!=c.conn{
                    fmt.Fprintf(c.conn, "\x1B[1G\x1B[2K%s has joined the chat\n[%s][%s]:", inUser.name, time.Now().Format("2006-01-02 15:04:05"), c.name)
				}
            }
        case msg := <-Messages:
            for _, c := range clients {
				frmt := fmt.Sprintf("[%s][%s]:%s", time.Now().Format("2006-01-02 15:04:05"), msg.sender.name, msg.text)
				mu.Lock()
				messages = append(messages, frmt)
				mu.Unlock()
				if msg.sender.conn != c.conn {
					fmt.Fprintf(c.conn, "\x1B[1G\x1B[2K%s[%s][%s]:", frmt, time.Now().Format("2006-01-02 15:04:05"), c.name)
				}
			}
        case outUser := <-quitMessages:
			for _, c := range clients {
				if outUser.conn!=c.conn{
                    fmt.Fprintf(c.conn, "\x1B[1G\x1B[2K%s has left the chat\n[%s][%s]:", outUser.name, time.Now().Format("2006-01-02 15:04:05"), c.name)
				}
			}	
		}
    }
}

func handleConn(conn net.Conn) {
	fmt.Fprintln(conn,"Welcome to TCP-Chat!")
	printPenguin(conn)
	var name string
	for {
		fmt.Fprint(conn, "[ENTER YOUR NAME]:")
		n, err := getName(conn)
		if err==nil {
			name = n
			break
		}
		fmt.Fprintf(conn, "%v\n", err)
	}
	
	user := User{
		name: name,
        conn: conn,
	}
	
	mu.Lock()
	clients = append(clients, user)
	mu.Unlock()

	for _, msg := range messages {
		fmt.Fprint(conn, msg)
	}

	welcomeMessages <- user

	for {
		fmt.Fprintf(conn, "[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), user.name )
		userInput, err :=  bufio.NewReader(user.conn).ReadString('\n')
		if err != nil {
			// fmt.Printf("Error reading from client: %s\n", err)
			break
		}
		if checkMsg(userInput) && checkLen(userInput){
			Messages <- Message{sender: user, text: userInput}
		}
	}

	defer func() {
		mu.Lock()
		for i, u := range clients {
			if u.conn == user.conn{
                clients = append(clients[:i], clients[i+1:]...)
                break
            }
		}
		mu.Unlock()
		quitMessages <- user
		user.conn.Close()
	}()  
	
}