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
		customerRouter.GET("check", serverApi.Check) // 创建客户
	}
}
