package server

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/MridulDhiman/dice/config"
	"github.com/MridulDhiman/dice/core"
)

// reading from the connection
func readFromConn(conn io.ReadWriter) (*core.RedisCmd, error) {
	// create a buffer of 512 bytes
	buf := make([]byte, 512);
	n, err:= conn.Read(buf) // it reads 512 bytes from connection and stores the bytes in bytes buffer, returning buffer size
	if err != nil {
		return nil, err
	}

	// decode the byte slice to array of string
	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, errors.New("cmd not found")
	}
	
	return &core.RedisCmd{
		Cmd: tokens[0],
		Args: tokens[1:],
	}, nil
}

func writeToConn(cmd *core.RedisCmd, conn io.ReadWriter) error {
	// writing to connection
	if err := core.EvalAndRespond(cmd, conn); err != nil {
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
		cmd, err:= readFromConn(conn)

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
			log.Println("command", cmd.Cmd)
			if err:= writeToConn(cmd, conn); err != nil {
				log.Println("err write: ", err)
			}

		}

	}
}