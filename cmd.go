package zerobot

import (
	"bytes"
	"os"
	"os/exec"
)

func (b *ZeroBot) handleSystemCmd(cmd string, args ...string) {
	if !isAllowCmds(cmd) {
		b.sendMsg("cmd:%s not allow", cmd)
		return
	}

	b.runCmdX(cmd, args...)
}

func (b *ZeroBot) runCmd(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

// run cmd, and send output/error to slack
func (b *ZeroBot) runCmdX(cmd string, args ...string) {
	buf := new(bytes.Buffer)
	c := exec.Command(cmd, args...)
	c.Stdout = buf
	c.Stderr = buf
	err := c.Run()
	if err != nil {
		b.sendErr(err)
		return
	}

	b.sendMsg(buf.String())
}

var allowCmds = map[string]struct{}{
	"top":  struct{}{},
	"ps":   struct{}{},
	"free": struct{}{},
	"df":   struct{}{},
	"htop": struct{}{},
	"ls":   struct{}{},
	"cat":  struct{}{},
	"du":   struct{}{},
	"env":  struct{}{},
}

func isAllowCmds(cmd string) bool {
	if _, ok := allowCmds[cmd]; ok {
		return true
	}
	return false
}
