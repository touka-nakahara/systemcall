package main

import (
	"path/filepath"
	"syscall"
)

func Cat(rp string) error {
	ap, err := filepath.Abs(rp)
	if err != nil {
		return nil
	}

	mfd, err := syscall.Open(ap, syscall.O_RDONLY, 666)

	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := syscall.Read(mfd, buf)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		// ほんとはキー入力待ちたい
	}

	syscall.Write(1, buf)

	return nil

}
