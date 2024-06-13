package fork

import (
	"fmt"
	"os"
	"syscall"
)

func Fork(fd int, rp string) {

	pa := &syscall.ProcAttr{
		Dir: "",
		Env: os.Environ(),
		Files: []uintptr{
			0,
			uintptr(fd), // ここに渡せばstdoutの書き込み先が変わる ( おもろい )
			2,
		},
		Sys: &syscall.SysProcAttr{}, // 誰？
	}
	pid, err := syscall.ForkExec(rp, []string{}, pa)

	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	data := fmt.Sprintf("Created Process %d from %d", pid, syscall.Getpid())
	syscall.Write(syscall.Stdout, []byte(data))

}
