package core

import (
	"errors"
	"fmt"
	"net"
)

func evalPING(args []string, conn net.Conn) error {
	// ping have only 1 argument that is, string and sends it back to the client
	if len(args) > 1 {
		return errors.New("ERR wrong no. of arguments for 'PING' command");
	}
	// no arguments provided
	var b []byte
	if len(args) == 0 {
		// encode output to simple/bulk string 
		b = Encode("PONG", true)
	} else {
		//else encode argument result to bulk string
		b = Encode(args[0], false)
	}

	// write to connection
	if _, err := conn.Write(b); err != nil {
		return err;
	}
	return nil
}

func EvalAndRespond(cmd *RedisCmd, conn net.Conn) error {
	fmt.Println("cmd: ", cmd.Cmd)

	if cmd.Cmd == "" {
		return errors.New("no cmd found");
	}

	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, conn)
	}
	return nil
}
