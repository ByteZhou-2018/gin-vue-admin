package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/constant"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type ServerApi struct{}

// 环境检测
func (s *ServerApi) Check(c *gin.Context) {

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Hour) // 设置为1小时
	defer cancel()

	// 将新的上下文传递给需要长时间运行的操作
	c.Request = c.Request.WithContext(ctx)

	checkDocker := cloud.CmdNode{
		Cmd: "docker version",
	}

	checkDockerCompose := cloud.CmdNode{
		Cmd: "docker-compose version",
	}

	Run([]cloud.CmdNode{checkDocker, checkDockerCompose}, c)

}

// 安装环境
func (s *ServerApi) Install(c *gin.Context) {

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Hour) // 设置为1小时
	defer cancel()

	// 将新的上下文传递给需要长时间运行的操作
	c.Request = c.Request.WithContext(ctx)

	docker := cloud.CmdNode{
		Cmd: "curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh",
	}

	dockerCompose := cloud.CmdNode{
		Cmd: "sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose && docker-compose version",
	}

	Run([]cloud.CmdNode{docker, dockerCompose}, c)

}

// 发布
func (s *ServerApi) Deploy() {

}

func Run(cmds []cloud.CmdNode, c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 将新的上下文传递给需要长时间运行的操作
	var checkRequest request.SSHRequest
	msg := make(chan cloud.MsgInfo, 100)
	err := c.ShouldBindJSON(&checkRequest)
	if err != nil {
		c.SSEvent(constant.FAIL, err.Error())
		return
	}

	client, err := checkRequest.Connection()
	if err != nil {
		c.SSEvent(constant.FAIL, err.Error())
		return
	}
	defer client.Close()

	go func() {
		c.Stream(func(w io.Writer) bool {
			select {
			case msg := <-msg:
				c.SSEvent(msg.Status, msg.Msg)
				if msg.Status == constant.COMPLETE {
					return false
				}
				return true
			case <-time.After(3 * time.Second):
				c.SSEvent(constant.PENDING, "pending:"+time.Now().String())
				return true
			}
		})
	}()

	for i := range cmds {
		cloudServiceGroup.Cmd(cmds[i], msg, client)
	}

	msg <- cloud.MsgInfo{Msg: "complete", Status: constant.COMPLETE}
}
