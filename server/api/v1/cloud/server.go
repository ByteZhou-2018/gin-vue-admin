package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"github.com/gin-gonic/gin"
	"io"
)

type ServerApi struct{}

// 环境检测
func (s *ServerApi) Check(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	var checkRequest request.SSHRequest
	msg := make(chan string, 1)
	err := c.ShouldBindJSON(&checkRequest)
	if err != nil {
		c.SSEvent("done", err.Error())
		return
	}

	ctx := context.Background()

	connect, err := checkRequest.Connection(ctx)

	if err != nil {
		c.SSEvent("done", err.Error())
		return
	}

	go cloudServiceGroup.Cmd("ping www.baidu.com", msg, *connect)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-msg:
			c.SSEvent("message", msg)
			return true
		case <-connect.CTX.Done():
			c.SSEvent("done", "")
			return false
		}
	})

}

// 安装环境
func (s *ServerApi) Install() {

}

// 发布
func (s *ServerApi) Deploy() {

}
