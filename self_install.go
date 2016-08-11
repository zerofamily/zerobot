package zerobot

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

const (
	binName = "zerobot"
	binDir  = "bin"
)

func (b *ZeroBot) handleBuild() {
	bin := getBinPath()

	md5Before := md5sum(bin)

	if err := b.runCmd("git", "pull", "-r"); err != nil {
		b.sendMsg("git pull err: %s", err.Error())
	}
	if err := b.runCmd("make", "build"); err != nil {
		b.sendMsg("make build err: %s", err.Error())
	}

	md5After := md5sum(bin)
	b.sendMsg("md5 before: %s\nmd5 after: %s", md5Before, md5After)
	if md5Before != md5After {
		b.sendMsg("built success")
	} else {
		b.sendMsg("nothing change")
	}
}

// send quit signal, and then supervisord let it restart automatically.
func (b *ZeroBot) handleRestart() {
	b.sendMsg("killing process ...")
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGHUP); err != nil {
		b.sendErr(err)
	}

	b.sendMsg("process killed")
}

func md5sum(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := md5.New()
	io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// bin path
func getBinPath() string {
	root, _ := os.Getwd()
	return filepath.Join(root, binDir, binName)
}
