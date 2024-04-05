package cloud

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ServerRouter struct{}

func (e *ServerRouter) InitServerRouter(Router *gin.RouterGroup) {
	customerRouter := Router.Group("server").Use(middleware.OperationRecord())
	serverApi := v1.ApiGroupApp.CloudApiGroup.ServerApi
	{
		customerRouter.POST("check", serverApi.Check)     // 检测环境
		customerRouter.POST("install", serverApi.Install) // 安装环境
		customerRouter.POST("zip", serverApi.Zip)         // 打包测试
		customerRouter.POST("deploy", serverApi.Deploy)   // 上传部署

	}
}
