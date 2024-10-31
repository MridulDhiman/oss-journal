## Abstract

PING command in redis is used to check whether the redis server is running or not, and if the connection is still alive.

Input: PING "hello world"
Output: "hello world"

Input: PING
Output: "PONG"

Only one argument: simple/bulk string

Usecases:
- Check whether the connection is still alive
- verifying the server's ability to serve data, by echoing back the same arguments, or "PONG" if no arguments
- measuring latency, time it takes for request to reach server.


## commit(a695284cd4eb4650286688574578b0c74059d92c): PING command implementation

### core/cmd.go: package core

- Create a struct `RedisCmd` for command(`string`) and it's arguments(`[]string`),  which is further used to evaluate PING command, when the RESP payload is deserialized.


### core/resp.go: 
#### DecodeArrayString()
- It decodes the input byte slice(received from Redis Client via RESP protocol) to any/interface{} type
- It then type asserts whether the output of `Decode()` call is of type `[]interface{}`, and if not, the program panics.
- it create a new array of string of tokens
- it type asserts each element of `[]interface` slice, whether it's of type string or not.
- if not, it panics the program
- return []string, error

#### Encode()
- if `isSimple == true`, encodes the value of type `interface{} or any` to simple string
- else, it gets encoded to byte string
- then we convert string to []byte slice and return the output

### server/sync_tcp.go: 
#### `readFromConn()`: returns *RedisCmd, error
- read 512 bytes data from connection
- stores it inside byte slice
- convert the byte slice to Array of strings using `DecodeArrayString()` function
- 1st element of array string will the command and rest of the elements will the command's arguments.
- create and return new struct of `RedisCmd` with the tokens of array strings.
- Command would in upper case 

> strings.ToUpper(str); // convert string to uppercase

### core/eval.go
#### `EvalAndRespond()`: returns error
- evaluating command name for "PING" from the switch case
- if no. of arguments > 1, return error
- if no. of arguments == 0, sends back "PONG" by encoding it to simple string, and converting it to byte slice
- else echos back the argument, by encoding it to bulk string and type converting it to byte slice
- write the byte slice to connection and if error return the error

## Conclusion
In this commit, the RESP encoded byte slice is being read from the connection, which is further decoded to array string to get commands and arguments. We write the reply back to the client, by evaluating the command ("PING" command in this case) and it's arguments.

In this commit, we evaluate the reply for "PING" command.  It takes a single argument,  which will return the same argument encoded to bulk string, if argument is provided. Otherwise we will return "PONG" encoded to simple string. Simple/bulk string is further converted to byte slice to write it to the TCP connection.

