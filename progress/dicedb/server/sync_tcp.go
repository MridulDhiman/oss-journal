package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/MridulDhiman/dice/config"
	"github.com/MridulDhiman/dice/core"
)

func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}
	return as, nil
}


// reading from the connection
func readFromConn(conn io.ReadWriter) ([]*core.RedisCmd, error) {
	// create a buffer of 512 bytes
	buf := make([]byte, 512)
	n, err := conn.Read(buf) // it reads 512 bytes from connection and stores the bytes in bytes buffer, returning buffer size
	if err != nil {
		return nil, err
	}

	values, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}

	var cmds []*core.RedisCmd = make([]*core.RedisCmd, 0)
	for _, value := range values.([]interface{}) {
		tokens, err := toArrayString(value.([]interface{}))
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}
	return cmds, nil

}

func writeToConn(cmd []*core.RedisCmd, conn io.ReadWriter) {
	// writing to connection
	core.EvalAndRespond(cmd, conn)

}

func RunSyncTCPServer() {
	lstnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	var conn_clients = 0
	if err != nil {
		panic(err)
	}

	for {
		// waiting for new client to connect: Blocking Call
		conn, err := lstnr.Accept() // and return new connection or error
		if err != nil {
			panic(err)
		}

		// once connection is established, we can start communicating

		// increment the no. of concurrent clients
		conn_clients += 1

		log.Println("client connected with address: ", conn.RemoteAddr(), "concurrent clients", conn_clients)
		cmds, err := readFromConn(conn)

		// could not read from connection: user disconnected
		if err != nil {
			// close the connection
			conn.Close()
			// decrement the no. of concurrent users
			conn_clients -= 1
			log.Println("client disconnected", conn.RemoteAddr(), "concurrent clients at this moment", conn_clients)

			// no more input left to be read: EOF(End Of File) Reached
			if err == io.EOF {
				break
			}

			// echo the data back to client
			writeToConn(cmds, conn)
		}

	}
}
