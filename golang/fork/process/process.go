package main

import (
	"fmt"
	"syscall"
)

func main() {

	for {
		pid := syscall.Getpid()
		syscall.Write(1, []byte(fmt.Sprintf("PID:%d", pid)))
	}

}
