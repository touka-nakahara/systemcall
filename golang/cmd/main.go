package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func main() {

	Cat("MEM")

	// STDOUTで出力する
	syscall.Write(syscall.Stdout, []byte("SYSCALL WRITE\n"))

	// プロセスID
	id := syscall.Getppid()
	syscall.Write(1, []byte(fmt.Sprintf("%d\n", id)))

	// SYSTEM CALLを使う
	message := []byte("SYSCALL SYSCALL6\n")
	syscall.Syscall6(syscall.SYS_WRITE, uintptr(1), uintptr(unsafe.Pointer(&message[0])), uintptr(len(message)), 0, 0, 0)

	// SYS_CREATはOPENのラッパー. Open(path, O_CREAT|O_WRONLY|O_TRUNC, mode らしい
	//DEAD うまくいかん
	// filePath := []byte("SYSCALL_CREAT")
	// _, _, errno := syscall.Syscall6(85, uintptr(unsafe.Pointer(&filePath[0])), 0666, 0, 0, 0, 0)
	// if errno != 0 {
	// 	fmt.Printf("SYSCALL ERROR : {%d}\n", errno)
	// }

	// Epolってなんだ？
	// ファイルを作る
	// MODE = syscall.O_CREAT|syscall.O_WRONLY|syscall.O_TRUNC
	// ファイルを作成する
	fd, _ := syscall.Open("SYSCALL_OPEN", syscall.O_CREAT|syscall.O_WRONLY|syscall.O_TRUNC, 0666)
	defer syscall.Close(fd)
	syscall.Write(fd, []byte(fmt.Sprintf("%d\n", id)))

	// /dev/nullのファイルディスクリプタ
	fdn, _ := syscall.Open("/dev/null", syscall.O_WRONLY, 0)
	syscall.Close(fdn)

	// REBOOT
	// 失敗した... カーネルランドで実行しないといけない？
	// sudoでやったら行けたけど, レポートが作成されたわ...
	_, _, errno := syscall.Syscall(syscall.SYS_REBOOT, 0x1234567, 0, 0)
	if errno != 0 {
		fmt.Printf("SYSCALL ERROR : {%s}\n", errno.Error())
	}

	//プロセスを作る

	// ファイルの中身を読む
	mfd, err := syscall.Open("MEM", syscall.O_RDONLY, 666)
	if err != nil {
		fmt.Println(err)
	}
	// pReadはオフセットあり？
	// サイズわかんなくね？?
	buf := make([]byte, 1024)
	for {
		n, err := syscall.Read(mfd, buf)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if n == 0 {
			break
		}
	}

	syscall.Write(1, buf)
	syscall.Write(1, []byte("\n"))

	// stdに吐き出す

	// Processにはこれを渡す
	// type ProcAttr struct {
	// 	Dir   string    // Current working directory.
	// 	Env   []string  // Environment.
	// 	Files []uintptr // File descriptors.
	// 	Sys   *SysProcAttr
	// }

	// メモリにつくる
	//DEAD できなかった
	// MemFdCreateしたかった...

	rp := "./build/process"
	syscall.SetNonblock(fd, true)

	// これって起動してるだけよな〜
	// プログラム化してプロセス起動したい
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
	fmt.Println(pid)
	fmt.Println(syscall.Getpid())

	for {
		// 排他処理になる...なんでだ
		syscall.Write(fd, []byte(fmt.Sprintf("PID:%d", syscall.Getpid())))
	}

	// STDIN

	// 仮想ファイル in memory のリソースでプロセスのリソース間をちょこちょこしたい

	// ForkLock だれ？

	// Creat, EpollCreate, EpollCreate1
	// Creatはファイルを作りそう

	// コマンド作ってみるのもおもろそうではあるね〜
	// CAT ~
	// ファイルの中身を標準出力〜

	// 	ts := syscall.Timespec{
	// 		Sec:  5,
	// 		Nsec: 0, // 誰？
	// 	}

	// 	syscall.Write(1, []byte("まつよ〜\n"))

	// 	// まってくれない...
	// 	_, _, errno = syscall.Syscall(35, uintptr(unsafe.Pointer(&ts)), uintptr(unsafe.Pointer(&ts)), 0)
	// 	if errno != 0 {
	// 		fmt.Printf("%s", errno.Error())
	// 	}

	//		syscall.Write(1, []byte("まった？\n"))
	//	}
}

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
	syscall.Write(1, []byte("\n"))

	return nil

}
