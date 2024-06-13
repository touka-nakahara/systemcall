package program

import (
	"fmt"
	"syscall"
)

func Hello() {
	id := syscall.Getpid()
	fmt.Printf("Hello World ! at %d from Printf!\n", id)
	syscall.Write(1, []byte(fmt.Sprintf("Hello World ! at %d from Syscall! \n", id)))
}
