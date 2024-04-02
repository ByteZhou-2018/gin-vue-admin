package cloud

import (
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
	Run("docker version && docker-compose version", c)
}

// 安装环境
func (s *ServerApi) Install(c *gin.Context) {
	docker := "curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh"

	dockerCompose := "sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose && docker-compose version"

	Run(docker+" ; "+dockerCompose, c)

}

// 发布
func (s *ServerApi) Deploy() {

}

func Run(cmd string, c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	var checkRequest request.SSHRequest
	done := make(chan int)
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
				if msg.Status == constant.FAIL {
					done <- 1
					return false
				}

				if msg.Status == constant.COMPLETE {
					done <- 1
					return false
				}
				return true
			case <-time.After(2 * time.Second):
				c.SSEvent(constant.PENDING, "pending")
				return true
			}
		})
	}()

	cloudServiceGroup.Cmd(cmd, msg, client)
	<-done
}
