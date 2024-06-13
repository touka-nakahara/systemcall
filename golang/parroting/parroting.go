package parroting

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	ShareKey   = 1234
	ShareSize  = 1024
	SemKey     = 5678
	IPC_CREATE = 01000
	IPC_RMID   = 0
)

var exit bool

func Parroting() {
	exit = false
	// 共有メモリへのアクセス
	// ない場合はIPC_CREATEで作成

	// 共有メモリの作成
	shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(ShareKey), uintptr(ShareSize), IPC_CREATE|0666)
	if errno != 0 {
		log.Fatalf("SYS_SHMGET: %s\n", errno.Error())
	}

	// 共有メモリのアタッチ
	shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, uintptr(shmid), 0, 0)
	if errno != 0 {
		log.Fatalf("SYS_SHMAT: %s\n", errno.Error())
	}
	defer syscall.Syscall(syscall.SYS_SHMDT, uintptr(shmaddr), 0, 0) // 共有メモリのデタッチ
	defer syscall.Syscall(syscall.SYS_SHMCTL, uintptr(shmid), 0, 0)  // 共有メモリの削除

	// 通知が来たら読み込んで返す
	semid, _, errno := syscall.Syscall(syscall.SYS_SEMGET, uintptr(SemKey), uintptr(1), IPC_CREATE|0666)
	if errno != 0 {
		log.Fatalf("SYS_SEMGET: %s\n", errno.Error())
	}

	// ロック解除を待つ ( 空きがあるならとる )
	type semBuf struct {
		sem_num uint16
		sem_op  int16
		sem_flg int16
	}
	sb := &semBuf{
		sem_num: 0,
		sem_op:  -1, // P操作
		sem_flg: 0,
	}
	_, _, errno = syscall.Syscall(syscall.SYS_SEMOP, uintptr(semid), uintptr(unsafe.Pointer(sb)), 0)
	if errno != 0 {
		log.Fatalf("SYS_SEMOP P Operation: %s\n", errno.Error())
	}

	// 読み込み
	rd := (*[ShareSize]byte)(unsafe.Pointer(shmaddr))[:]
	syscall.Write(1, rd)
	syscall.Write(1, []byte("\n"))

	// 書き込み
	pid := syscall.Getppid()
	data := fmt.Sprintf("Hello I'm at %d from stdout", pid)
	syscall.Write(1, []byte(data))
	copy((*[ShareSize]byte)(unsafe.Pointer(shmaddr))[:], data)

	// 書き込める

	// for !exit {
	// 	// SIGINTが来たらBreakして終了すれば良い
	// 	data := fmt.Sprintf("Hello Sharing Memory I'm at %d", pid)
	// 	copy((*[ShareSize]byte)(unsafe.Pointer(shmaddr))[:], data)
	// }

	for {

	}
}
