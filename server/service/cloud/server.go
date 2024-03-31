package cloud

import (
	"bufio"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
)

type ServerService struct{}

// 环境检测
func (s *ServerService) Cmd(cmd string, msg chan string, connect cloud.Context) error {
	client := connect.Client
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)

	err = session.Start(cmd)

	if err != nil {
		return err
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			connect.Cancel()
			break
		}
		msg <- line
	}
	return nil
}
