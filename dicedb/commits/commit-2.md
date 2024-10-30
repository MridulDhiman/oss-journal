## Abstract 

RESP(Redis Serialization Protocol): 
Redis client uses RESP protocol to communicate with Redis Server.

Advantages of RESP protocol:
- simple to implement
- fast to parse 
- human readable

It can serialize almost any data type like text, number, string. Client sends commands in form of array of bulk strings(bulk string is a binary string of 512MB size). The client sends command and it's arguments. 1st string of the array (or second sometimes), contains the command name and rest of the array is allocated for the command's arguments. RESP Server parses the command, does processing on basis of command and returns result according to command's data types.
It is not TCP specific, but generally uses TCP in the network layer, establishing TCP 3-way connection before any message being sent. TCP ensures that the data is not lost via sequence numbers and congestion control between the routers is managed properly.


All the data is packaged inside a CRLF strings, which is abbreviation for Carriage Return (`\r`), Line Feed (`\n`)
E.g. if Redis Client wants to send "hello" message to Redis Server, this is how the actual message payload is formatted: 

`$5\r\nhello\r\n`.

`hello` contains 5 characters and each character as we know is 1 byte each, thus string would of 5 bytes in total. We prefix the length in the start of the payload, which ensures that we don't have to keep listening to the process for any end marker (like EOF etc.).
In contrast, other client-server communication protocols like HTTP have a `newline` marker signifying the end of the payload.
Technically, both of them are different as in, one is "You have 100 bytes of data to read", while other "You have to keep listening for new data, until you read the word STOP". First one is more intuitive and reliable.


RESP protocol is considered binary safe, as it can serialize any data type like text, string, numbers, arrays. It also has support for error type serialization.
It can also serialize raw binary types like images, files etc. Some protocols lack in recognising byte patterns and binary sequences, but that is not the case with RESP protocol.


It was introduced in Redis 1.2 version as in RESP1 specification. Today, redis also has RESP2 and RESP3 specifications.
We are going to implement RESP1 specification in this commit.

## Commit(7f3265f5e941d5cf2cb95de022f3c3b46f79937d) RESP1 Specification Decoder and it's tests

### `core/resp.go` : package `core`
- Data is in form of bytes buffer
- first byte is prefix notation for data type, rest bytes contains type's content.

#### Protocol Format: RESP uses prefix notation where each data type is marked with a special character:
1. `+` : for simple string
2. `-` : for errors
3. `:` : for integers
4. `$` : for bulk strings(binary string with CRLF in between)
5. `*` : for arrays
---
- each line ends with CRLF`\r\n`.
- Each data type returns data, delta(no. of bytes consumed) and error

#### For Simple string
- implements `readSimpleString()` function
- reads RESP encoded simple string from data and returns string, delta and error
- Traverse till `\r` reached
- delta(no. of bytes consumed): len + 2(\r\n)

#### for error
- reads RESP encoded error from data and return error string, delta and error
- uses `readSimpleString()` function

#### for 64 bit integer 
- reads RESP encoded integer from data and return integer, delta and error
- traverse till `\r` reached
- to get integer value, traverse each character and convert to number and add it to result*10 to get desired result.

`result = result*10 + int64(data[pos] - '0');`

#### For bulk strings
- in case of bulk strings, length of the strings is also prefixed before actual data
- it is in this format: `$5\r\nhello\r\n`

#### for arrays
- It also has prefixed length
- create new array of `any` or `interface{}` type with that length
- traverse the array and decode the elements and increment position using delta

---
### `core/resp_test.go`:  package `core_test`

In golang, there are 2 ways of naming tests, either keeping tests in the same package, suitable for white box testing, where I can check out the inner implementation of the packages or using different package names, suitable for blackbox testing, where we check whether for particular input, we get the desired output or not.

- we use built-in `testing` library for unit tests
- here we write unit tests which are black box in nature.

#### Simple String Tests
- Input: `+OK\r\n`
- Expected Output: `OK`

#### Error Tests
- Input: `-Error Message\r\n`
- Expected Output: `Error Message`

#### Int64 Tests
- Input: `:1000\r\n`
- Expected Output: 1000

#### Bulk String Tests
- Input: `$5\r\nhello\r\n`
- Output: "hello"

#### array tests
- Input: `*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n`
- Expected Output: `{"hello", "world"}`
- Here length specify the length of the array
---
- Length is written as `$5\r\n` or `*5\r\n`
















