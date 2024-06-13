package main

import (
	"fmt"
	"syscall"
	"syscall/fork"
)

func main() {
	rp := "./fork/build/process" // processの場所

	fd, _ := syscall.Open("SYSCALL_OPEN", syscall.O_CREAT|syscall.O_WRONLY|syscall.O_TRUNC, 0666)
	defer syscall.Close(fd)

	fork.Fork(fd, rp)

	for {
		// 排他処理になる...なんでだ
		syscall.Write(fd, []byte(fmt.Sprintf("PID:%d", syscall.Getpid())))
	}

}
