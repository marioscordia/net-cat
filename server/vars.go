package server

import (
	"net"
	"sync"
)

type Message struct {
	sender User
	text string
}

var ( 
	mu sync.Mutex
	clients  []User
	messages []string
	pinguin = []string{
	"         _nnnn_",
	"        dGGGGMMb",
	"       @p~qp~~qMb",
	"       M|@||@) M|",
	"       @,----.JM|",
	"      JS^\\__/  qKL",
	"     dZP        qKRb",
	"    dZP          qKKb",
	"   fZP            SMMb",
	"   HZM            MMMM",
	"   FqM            MMMM",
	" __| \".        |\\dS\"qML",
	" |    `.       | `' \\Zq",
	"_)      \\.___.,|     .'",
	"\\____   )MMMMMP|   .'",
	"     `-'       `--'",
}
)

type User struct {
	name string
	conn net.Conn
}