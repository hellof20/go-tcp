package main
 
import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	proxyproto "github.com/pires/go-proxyproto"
)
 
func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	proxyListener := &proxyproto.Listener{Listener: listener}
	defer proxyListener.Close()
	
	for {
		con, err := proxyListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClientRequest(con)
	}
}
 
func handleClientRequest(con net.Conn) {
	addr := con.RemoteAddr()
	defer con.Close()
	clientReader := bufio.NewReader(con)
	for {
		clientRequest, err := clientReader.ReadString('\n')
		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println(addr,"client requested to close")
				return
			} else {
				log.Println(addr,clientRequest)
			}
		case io.EOF:
			// log.Println(addr,"client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
 
		// Responding to the client request
		if _, err = con.Write([]byte("GOT IT!\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}
