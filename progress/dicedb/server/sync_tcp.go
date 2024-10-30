package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/MridulDhiman/dice/config"
)

// reading from the connection
func readFromConn(conn net.Conn) (string, error) {
	// create a buffer of 512 bytes
	buf := make([]byte, 512);
	n, err:= conn.Read(buf) // it reads 512 bytes from connection and stores the bytes in bytes buffer, returning buffer size
	if err != nil {
		return "", err
	}
	
	return string(buf[:n]), nil
}

func writeToConn(data string, conn net.Conn) error {
	// writing to connection
	if _, err := conn.Write([]byte(data)); err != nil {
		return err
	}
	return nil
}

func RunSyncTCPServer() {
	lstnr, err := net.Listen("tcp", config.Host+ ":" + strconv.Itoa(config.Port));
	var conn_clients = 0
	if err != nil {
		panic(err)
	}

	for {
		// waiting for new client to connect: Blocking Call
		conn, err:= lstnr.Accept() // and return new connection or error
		if err != nil {
			panic(err)
		}

		// once connection is established, we can start communicating

		// increment the no. of concurrent clients
		conn_clients+= 1;

		log.Println("client connected with address: ", conn.RemoteAddr(), "concurrent clients", conn_clients)
		data, err:= readFromConn(conn)

		// could not read from connection: user disconnected
		if err != nil {
             // close the connection
			 conn.Close()
			// decrement the no. of concurrent users
			conn_clients-=1;
			log.Println("client disconnected", conn.RemoteAddr(), "concurrent clients at this moment", conn_clients)
			
			// no more input left to be read: EOF(End Of File) Reached
			if err == io.EOF {
				break
			}

			// echo the data back to client
			log.Println("command", data)
			if err:= writeToConn(data, conn); err != nil {
				log.Println("err write: ", err)
			}

		}

	}
}