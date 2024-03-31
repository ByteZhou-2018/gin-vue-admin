package cloud

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type ServerApi struct{}

// 环境检测
func (s *ServerApi) Check(c *gin.Context) {
	var checkRequest request.SSHRequest
	msg := make(chan string, 1)
	err := c.ShouldBindJSON(&checkRequest)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	ctx := context.Background()

	connect, err := checkRequest.Connection(ctx)

	go cloudServiceGroup.Cmd("docker-compose version", msg, *connect)

	var remsg string

	for {
		select {
		case msg := <-msg:
			remsg += msg
		case <-connect.CTX.Done():
			response.OkWithMessage(remsg, c)
			return
		default:
			break
		}
	}
}

// 安装环境
func (s *ServerApi) Install() {

}

// 发布
func (s *ServerApi) Deploy() {

}
