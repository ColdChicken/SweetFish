package common

import (
	xe "be/common/error"
	"be/common/log"
	"bytes"
	"os/exec"
	"strings"
	"time"
)

// Exec 执行命令
func Exec(timeout int64, cmd string, args ...string) (string, string, error) {
	log.Debugf("执行命令 %s %s", cmd, strings.Join(args, " "))
	command := exec.Command(cmd, args...)
	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	command.Stdout = &cmdOut
	command.Stderr = &cmdErr

	err := command.Start()
	if err != nil {
		log.Errorf("执行 %s 命令失败,  %s", cmd, err.Error())
		return "", "", err
	}

	done := make(chan error)
	go func() {
		done <- command.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		log.Errorf("执行 %s 命令超时，准备kill进程", cmd)
		if err := command.Process.Kill(); err != nil {
			log.Errorf("执行 %s 命令超时，kill进程失败", cmd)
		}
		go func() {
			<-done
		}()
		return "", "", xe.New("执行超时")
	case err = <-done:
		if err != nil {
			log.Errorf("执行 %s 异常, %s, stdout %s, stderr %s", cmd, err.Error(), cmdOut.String(), cmdErr.String())
			return "", "", err
		} else {
			log.Debugf("执行 %s 成功", cmd)
			return cmdOut.String(), cmdErr.String(), nil
		}
	}
}
