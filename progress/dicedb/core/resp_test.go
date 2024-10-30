package core_test

import (
	"fmt"
	"testing"

	"github.com/MridulDhiman/dice/core"
)

func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n" :"OK",
	}

	for k,v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			// test fails
			t.Logf("Expected %v, Got %v", v, value)
			t.Fail()
		}
	}
}

func TestErrorDecode(t *testing.T) {
	cases := map[string]string {
		"-Error Message\r\n" : "Error Message",
	}
	for k,v  := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Logf("Expected %v, Got %v", v, value)
			t.Fail()
		}
	}
}

func TestInt64Decode(t *testing.T) {
	cases := map[string]int64 {
		":0r\n": 0,
		":1000\r\n" : 1000,
	}
	for k,v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Logf("Expected %v, Got %v", v, value)
			t.Fail()
		}
	}
}

func TestBulkStringDecode (t *testing.T) {
	cases := map[string]string {
		"$5\r\nhello\r\n": "hello",
		"$0\r\n\r\n": "",
	}

	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		if value != v {
			t.Logf("Expected %v, Got %v", v, value)
			t.Fail()
		}
	}
}


func TestArrayDecode(t *testing.T) {
	cases:= map[string][]interface{} {
		"*0\r\n" : {},
		"*2\r\n+hello\r\n+world\r\n" : {"hello", "world"},
	}
	for k, v := range cases {
		value, _ := core.Decode([]byte(k))
		// type assertion syntax to convert interface to slice of interface
		array, ok:= value.([]interface{})
		if !ok {
			t.Log("could not convert interface to slice of interface")
			t.Fail()
		}
		if len(array) != len(v) {
			t.Log("Length of slice does not match")
			t.Logf("Expected %v, Got %v", len(v), len(array))
			t.Fail()
		}

		for i := range array {
			// this can be done, but it will cause test to fail if type does not matches
			// if v[i] != array[i] {
			// 	t.Fail()
			// }
			// instead normalize all the types to string first, and then compare
			if fmt.Sprintf("%v", v[i]) != fmt.Sprintf("%v", array[i]) {
				t.Fail()
			}
		}
	}
}