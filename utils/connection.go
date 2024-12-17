package utils

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	port  = "8989"
	mu    sync.Mutex
	USERS = make(map[string]net.Conn)
)

func HandleConnection(conn *net.Conn, PORT string) {
	port = PORT
	if len(USERS) == 10 {
		(*conn).Write([]byte("Room is full only 10 people allowed"))
		(*conn).Close()
		return
	}
	greeting := "Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'"
	defer (*conn).Close()
	(*conn).Write([]byte(greeting))
	name := login(conn, 0)
	if name == "" {
		(*conn).Write([]byte("\033[F\ntoo many attempts"))

		(*conn).Close()
		return
	}
	chat(conn, &name)
	disconect(conn, name)
	delete(USERS, name)
}

func login(conn *net.Conn, spam int) string {
	if spam == 5 {
		return ""
	}
	connFile, err := os.OpenFile("netcat-connection_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return ""
	}
	defer connFile.Close()
	date := time.Now().Format(time.DateTime)
	username := ""
	buffer := make([]byte, 1024)
	(*conn).Write([]byte("\n[ENTER YOUR NAME]:"))
	for {
		n, err := (*conn).Read(buffer)
		if err == io.EOF {
			return ""
		}
		username += string(buffer[:n-1])
		if err != nil {
			break
		}
		if strings.Contains(string(buffer), "\n") {
			break
		}
	}
	status := checkUsername(username, conn)
	if status != "" {
		(*conn).Write([]byte(status))
		return login(conn, spam+1)
	} else {
		oldchat, err := os.ReadFile("netcat-chat_" + port + ".log")
		if err != nil {
			(*conn).Write([]byte("connot access oldchat\n[" + date + "][" + username + "]:"))
		} else {
			(*conn).Write(oldchat)
			(*conn).Write([]byte("[" + date + "][" + username + "]:"))
		}
	}
	for user, Conn := range USERS {
		if user != username {
			Conn.Write([]byte("\n" + username + " has joined our chat...\n[" + date + "][" + user + "]:"))
		}
	}

	connFile.Write([]byte(username + " has joined our chat...\n"))
	return username
}

func checkUsername(username string, conn *net.Conn) string {
	mu.Lock()
	defer mu.Unlock()

	if len(username) < 3 {
		return "username too small"
	}
	if USERS[username] != nil {
		return "username already used"
	}
	if len(username) > 25 {
		return "username too long"
	}
	if !validchars(username) {
		return "only use latin letters and \"-\""
	}
	if len(USERS) == 10 {
		return "room is full"
	}
	USERS[username] = (*conn)
	return ""
}

func validchars(s string) bool {
	for _, v := range s {
		if !((v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || v == '-') {
			return false
		}
	}
	return true
}

func chat(Conn *net.Conn, name *string) {
	chatFile, err := os.OpenFile("netcat-chat_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln(err)
	}
	msgPrefix := "[" + time.Now().Format(time.DateTime) + "][" + (*name) + "]:"
	msg := ""
	buffer := make([]byte, 1024)
	for {
		n, err := (*Conn).Read(buffer)
		if err != nil {
			return
		}
		msg += string(buffer[:n-1])
		if strings.Contains(string(buffer), "\n") {
			break
		}
	}
	msg = strings.TrimSpace(msg)
	if !Validmsg(msg, Conn) {
		(*Conn).Write([]byte("\033[2K[" + time.Now().Format(time.DateTime) + "][" + (*name) + "]:"))
	} else {
		for name, conn := range USERS {
			if conn != (*Conn) {
				conn.Write([]byte("\a\n" + msgPrefix + msg + "\n"))
			}
			conn.Write([]byte("[" + time.Now().Format(time.DateTime) + "][" + name + "]:"))

		}
		fmt.Fprint(chatFile, msgPrefix+msg+"\n")
	}
	chatFile.Close()
	chat(Conn, name)
}

func disconect(conn *net.Conn, name string) {
	connFile, err := os.OpenFile("netcat-connection_"+port+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer connFile.Close()
	for user, c := range USERS {
		if c != (*conn) {
			c.Write([]byte("\n" + name + " has left our chat...\n[" + time.Now().Format(time.DateTime) + "][" + user + "]"))
		}
	}
	connFile.Write([]byte(name + " has left our chat...\n"))
}

func Validmsg(msg string, conn *net.Conn) bool {
	if len(msg) > 255 {
		fmt.Fprintln((*conn), "message too long....")
		return false
	}
	if msg == "" {
		fmt.Fprint((*conn), "\033[F")
		return false
	}
	for _, v := range msg {
		if (v < 32 || v > 126) && (v < 128 || v > 255) {
			fmt.Fprintln((*conn), "invalid characters....")
			return false
		}
	}
	return true
}
