package core

import "errors"


// reads the RESP encoded simple string from the data and returns
// string, delta(no. of bytes consumed) and error
func readSimpleString(data []byte) (string, int, error) {
	var pos int = 1
	for ; data[pos] != '\r'; pos++ {}
	// 1 based
	delta := pos + 2
	return string(data[1:pos]), delta, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInteger(data []byte) (int64, int, error) {
	pos := 1
	var result int64 = 0
	for ; data[pos] != '\r'; pos++ {
		result = result * 10 + int64(data[pos]-'0');
	}
	return result, pos + 2, nil
}


func readLength(data []byte) (int, int) {
	pos, length := 0,0
	for pos = range data {
		b:= data[pos]
		// break when /r reached
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length * 10 + int(b - '0')
	}
	return 0,0
}

func readBulkString(data []byte) (string, int, error) {
	pos:= 1
	len, delta := readLength(data[pos:])
	// pos will be at actual bulk string
	pos += delta
	bulkStr := string(data[pos:(pos+len)])
	return bulkStr, pos + len + 2, nil
}

func readArray(data []byte) (interface{}, int, error) {
	pos:=1
	len, delta:= readLength(data[pos:])
	pos+= delta

	// create an slice of any/interface{} type
	var elems []interface{} = make([]interface{}, len)
   // index
	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, errors.New("could not decode array element")
		}
		elems[i] = elem
		// bytes consumed
		pos += delta
	}
	// all the bytes would be consumed 
	return elems, pos, nil
}

// 1st byte signifies data type 
func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data found")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data);
	case '-':
		return readError(data);
	case ':':
		return readInteger(data);
	case '$':
		return readBulkString(data);
	case '*':
		return readArray(data);
	}
	return nil, 0, nil;
}

// data in form of bytes buffer
func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data found")
	}
	value, _, err := DecodeOne(data)
	return value, err
}