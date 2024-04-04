package cloud

import (
	"bufio"
	"github.com/flipped-aurora/gin-vue-admin/server/constant"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"golang.org/x/crypto/ssh"
)

type ServerService struct{}

// 环境检测
func (s *ServerService) Cmd(cmd cloud.CmdNode, msg chan cloud.MsgInfo, client *ssh.Client) {

	session, err := client.NewSession()
	if err != nil {
		msg <- cloud.MsgInfo{Msg: err.Error(), Status: constant.FAIL}
		return
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		msg <- cloud.MsgInfo{Msg: err.Error(), Status: constant.FAIL}
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		msg <- cloud.MsgInfo{Msg: err.Error(), Status: constant.FAIL}
		return
	}

	// 创建Scanner来逐行读取流
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	// 启动协程读取stdout
	go func() {
		for stdoutScanner.Scan() {
			msg <- cloud.MsgInfo{Msg: stdoutScanner.Text(), Status: constant.MESSAGE}
		}
	}()

	// 启动协程读取stderr
	go func() {
		for stderrScanner.Scan() {
			msg <- cloud.MsgInfo{Msg: stderrScanner.Text(), Status: constant.FAIL}
		}
	}()

	err = session.Start(cmd.Cmd)

	if err != nil {
		msg <- cloud.MsgInfo{Msg: err.Error(), Status: constant.FAIL}
	}

	msg <- cloud.MsgInfo{Msg: "run:" + cmd.Cmd, Status: constant.MESSAGE}

	err = session.Wait()
	if err != nil {
		msg <- cloud.MsgInfo{Msg: err.Error(), Status: constant.FAIL}
	}

}
