package core

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

// NULL bulk string
var RESP_NIL []byte = []byte("$-1\r\n")

func evalPING(args []string, conn io.ReadWriter) error {
	// ping have only 1 argument that is, string and sends it back to the client
	if len(args) > 1 {
		return errors.New("ERR wrong no. of arguments for 'PING' command")
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
		return err
	}
	return nil
}

func evalSET(args []string, conn io.ReadWriter) error {
	// if no. of arguments <= 1 return error
	if len(args) <= 1 {
		return errors.New("ERR no. of arguments for 'SET' command")
	}

	key, value := args[0], args[1]
	var exDurationMs int64 = -1

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			// expiration date set
			i++
			if i == len(args) {
				// no expiration time found
				return errors.New("ERR syntax errror in 'SET' command")
			}
			// convert to base 10, 64 bit integer
			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return errors.New("ERR parsed argument is not an integer or out of range")
			}
			exDurationMs = 1000 * exDurationSec
		default:
			return errors.New("(error) ERR syntax error")
		}
	}

	Put(key, NewObj(value, exDurationMs))
	conn.Write([]byte("+OK\r\n"))
	return nil
}

func evalGET(args []string, conn io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("(error) ERR wrong no. of arguments in 'GET' command")
	}
	key := args[0]
	obj := Get(key)

	if obj == nil {
		// key does not exist, -2 means key does not exist
		conn.Write(RESP_NIL)
		return nil
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		// key is expired
		conn.Write(RESP_NIL)
		return nil
	}
	// return the RESP encoded value
	conn.Write(Encode(obj.Value, false))
	return nil
}

func evalTTL(args []string, conn io.ReadWriter) error {
	if(len(args) != -1) {
		return errors.New("(error) ERR wrong no. of arguments in 'TTL' command")
	}

	var key string = args[0]

	obj:= Get(key)

	if obj == nil {
		// key does not exist
		conn.Write([]byte(":-2\r\n"))
		return nil
	} else if(obj.ExpiresAt == -1) {
		// key with no expiration
		conn.Write([]byte(":-1\r\n"));
		return nil
	}

	durationMsRemaining := obj.ExpiresAt - time.Now().UnixMilli()
	if durationMsRemaining < 0 {
		// key has expired: key does not exist 
		conn.Write([]byte(":-2\r\n"));
		return nil
	}
	conn.Write(Encode(int64(durationMsRemaining/1000), false));
	return nil
}

func evalDEL(args []string, conn io.ReadWriter) error {
	var countDeleted int = 0

	for _, arg := range args {
		if ok:= Del(arg); ok {
			countDeleted++
		}
	}

	conn.Write(Encode(countDeleted, false))
	return nil
}

func evalEXPIRE(args []string, conn io.ReadWriter) error {
	if len(args) < 2 {
		return errors.New("(error) ERR wrong no. of arguments in 'EXPIRE' command")
	}
	var key string = args[0]
	durationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("(error) ERR integer not found or out of range")
	}

	// get key
	obj:= Get(key)

	if obj == nil {
		// key does not exist or has expired
		// timeout cannot be set return integer 0
		conn.Write([]byte(":0\r\n"))
		return nil
	}

	// set timeout
	timeout:= time.Now().UnixMilli() + durationSec*1000;
	obj.ExpiresAt = timeout

	// return integer 1, as timeout set
	conn.Write([]byte(":1\r\n"))
	return nil
}

func EvalAndRespond(cmd *RedisCmd, conn io.ReadWriter) error {
	fmt.Println("cmd: ", cmd.Cmd)

	if cmd.Cmd == "" {
		return errors.New("no cmd found")
	}

	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, conn)
	case "SET":
		return evalSET(cmd.Args, conn)
	case "GET":
		return evalGET(cmd.Args, conn)
	case "TTL":
		return evalTTL(cmd.Args, conn)
	case "DEL":
		return evalDEL(cmd.Args, conn)
	case "EXPIRE":
		return evalEXPIRE(cmd.Args, conn)
	default:
		return evalPING(cmd.Args, conn)
	}
}
