package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://blog.golang.org/errors-are-values
	b := bufio.NewWriter(os.Stdout)
	b.Write([]byte("hi, "))
	b.Write([]byte("there"))
	b.Write([]byte("\n"))
	if b.Flush() != nil {
		fmt.Println(b.Flush().Error())
	}
}
