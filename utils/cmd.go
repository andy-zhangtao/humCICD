/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/28.

// CmdRun 封装好的命令执行函数
// 执行结束后,统一返回命令的标准输出和错误输出
func CmdRun(command []string) (string, error) {
	var out, stderr bytes.Buffer
	var cmd *exec.Cmd

	cmd = exec.Command(command[0], command[1:]...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	return fmt.Sprintf("%s\n%s", out.String(), stderr.String()), err
}
