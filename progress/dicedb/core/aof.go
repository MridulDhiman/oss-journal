package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MridulDhiman/dice/config"
)

func DumpAllKeys() {
	// open file in write only file and creates new file if not created
	file, err:= os.OpenFile(config.AOFFile, os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
	}

// rewriting AOF from scratch 
	for k, v := range store {
		cmd := fmt.Sprintf("SET %s %s", k, v.Value)
		tokens := strings.Split(cmd, " ")
		file.Write(Encode(tokens, false))
	}

	fmt.Println("AOF rewrite successful")
}