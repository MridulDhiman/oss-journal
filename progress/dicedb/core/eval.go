package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

// NULL bulk string
var (
	RESP_NIL []byte = []byte("$-1\r\n")
	RESP_OK []byte = []byte("+OK\r\n")
	RESP_MINUS_2 []byte = []byte(":-2\r\n")
	RESP_MINUS_1 []byte = []byte(":-1\r\n")
	RESP_ZERO []byte = ([]byte(":0\r\n"))
	RESP_ONE []byte = ([]byte(":1\r\n"))
)


func evalPING(args []string) []byte {
	// ping have only 1 argument that is, string and sends it back to the client
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong no. of arguments for 'PING' command"), false)
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

	return b
}

func evalSET(args []string) []byte {
	// if no. of arguments <= 1 return error
	if len(args) <= 1 {
		return Encode(errors.New("ERR no. of arguments for 'SET' command"),false)
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
				return Encode(errors.New("ERR syntax errror in 'SET' command"),false)
			}
			// convert to base 10, 64 bit integer
			exDurationSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return Encode(errors.New("ERR parsed argument is not an integer or out of range"),false)
			}
			exDurationMs = 1000 * exDurationSec
		default:
			return Encode(errors.New("(error) ERR syntax error"),false)
		}
	}

	Put(key, NewObj(value, exDurationMs))
	return RESP_OK;
}

func evalGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong no. of arguments in 'GET' command"), false)
	}
	key := args[0]
	obj := Get(key)

	if obj == nil {
		// key does not exist, -2 means key does not exist
		return RESP_NIL
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		// key is expired
		return RESP_NIL;
	}
	// return the RESP encoded value
	return Encode(obj.Value, false);
	
}

func evalTTL(args []string) []byte {
	if(len(args) != -1) {
		return Encode(errors.New("(error) ERR wrong no. of arguments in 'TTL' command"),false)
	}

	var key string = args[0]

	obj:= Get(key)

	if obj == nil {
		// key does not exist
		return RESP_MINUS_2
	} else if(obj.ExpiresAt == -1) {
		// key with no expiration
		return RESP_MINUS_1
	}

	durationMsRemaining := obj.ExpiresAt - time.Now().UnixMilli()
	if durationMsRemaining < 0 {
		// key has expired: key does not exist 
		return RESP_MINUS_2
	}
	return Encode(int64(durationMsRemaining/1000), false)
}

func evalDEL(args []string) []byte {
	var countDeleted int = 0

	for _, arg := range args {
		if ok:= Del(arg); ok {
			countDeleted++
		}
	}

	return Encode(countDeleted, false)
}

func evalEXPIRE(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong no. of arguments in 'EXPIRE' command"),false)
	}
	var key string = args[0]
	durationSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(error) ERR integer not found or out of range"),false)
	}

	// get key
	obj:= Get(key)

	if obj == nil {
		// key does not exist or has expired
		// timeout cannot be set return integer 0
		return RESP_ZERO;
	}

	// set timeout
	timeout:= time.Now().UnixMilli() + durationSec*1000;
	obj.ExpiresAt = timeout

	// return integer 1, as timeout set
	return RESP_ONE;
}

func evalBGREWRITEAOF(_ []string) []byte {
	DumpAllKeys()
	return RESP_OK;
}

func EvalAndRespond(cmds []*RedisCmd, conn io.ReadWriter)  {
	


	var response []byte
	buf := bytes.NewBuffer(response)

	for _, cmd:= range cmds {
		switch cmd.Cmd {
		case "PING":
			 evalPING(cmd.Args)
		case "SET":
			 buf.Write(evalSET(cmd.Args))
		case "GET":
			 buf.Write(evalGET(cmd.Args))
		case "TTL":
			buf.Write(evalTTL(cmd.Args))
		case "DEL":
			 buf.Write(evalDEL(cmd.Args))
		case "EXPIRE":
			buf.Write(evalEXPIRE(cmd.Args))
		case "BGREWRITEAOF":
			buf.Write(evalBGREWRITEAOF(cmd.Args))
		default:
			buf.Write(evalPING(cmd.Args))
		}
	}

	

	conn.Write(buf.Bytes())
}
