package main

import "syscall"

func main() {
	b := make([]byte, 1024)
	// 入力待ちになる？
	syscall.Read(0, b)

	syscall.Write(1, b)
}
