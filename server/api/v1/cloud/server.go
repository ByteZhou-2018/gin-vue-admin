package cloud

import (
	"archive/zip"
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/constant"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ServerApi struct{}

// 整体为使用docker-compose部署项目
// 前端使用nginx容器,nginx.conf需要映射到本地文件
// 打包待发布的项目，先找到前端目录，在目录下执行npm run build，获得dist目录
// 然后到后端文件下，根据传入的不同的system系统，设置环境变量，执行go build -o server
// windows变量 GOOS=windows GOARCH=amd64
// linux变量 GOOS=linux GOARCH=amd64
// mac变量 GOOS=darwin GOARCH=amd64
// 然后将前端dist目录和后端server文件+ 后端resource + 后端 config.yaml 打包成zip文件
// 根据config.yaml文件中的system的数据库类型设置和数据库具体配置来创建docker-compose.yml文件
// 如果数据库地址不为127.0.0.1，则需要在docker-compose.yml文件中设置数据库地址
// 如果数据库地址为127.0.0.1，则需要创建mysql容器，并且在docker-compose.yml文件中设置数据库地址，且挂盘数据到本地
// 然后将zip文件传输到服务器上，解压到服务器/www-gva目录下
// 进入/www-gva目录下，执行docker-compose up -d
func (s *ServerApi) Zip(c *gin.Context) {
	var de request.Deploy
	err := c.ShouldBindJSON(&de)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	webPath := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Web, "..")

	// 执行一个命令，获取错误输出
	webCmd := exec.Command("npm", "run", "build")
	webCmd.Dir = webPath
	webOut, err := webCmd.CombinedOutput()
	if err != nil {
		response.FailWithMessage(string(webOut), c)
		return
	}

	serverCmd := exec.Command("go", "build", "-o", "gin-vue-admin")
	switch de.SystemType {
	case "windows":
		serverCmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64")
	case "linux":
		serverCmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	case "mac":
		serverCmd.Env = append(os.Environ(), "GOOS=darwin", "GOARCH=amd64")
	}

	serverPath := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server)
	serverCmd.Dir = serverPath
	serverOut, err := serverCmd.CombinedOutput()

	if err != nil {
		response.FailWithMessage(string(serverOut), c)
		return
	}

	fileName := "deploy.zip"
	// 创建一个新的zip文件
	zipFile, err := os.Create(fileName)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer zipFile.Close()

	// 创建一个zip写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = utils.DoZip(zipWriter, filepath.Join(webPath, "dist"), "dist")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("打包成功", c)
}

func (s *ServerApi) Down(c *gin.Context) {
	//Run([]cloud.CmdNode{checkDocker, checkDockerCompose}, c)
}

func (s *ServerApi) Deploy(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 将新的上下文传递给需要长时间运行的操作
	var checkRequest request.SSHRequest
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

	serverPath := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server)

	serverFile := filepath.Join(serverPath, "gin-vue-admin")
	webFile := filepath.Join(serverPath, "deploy.zip")

	err = SftpUpload(serverFile, "/www-gva/server", "/www-gva/server/gin-vue-admin", client)
	if err != nil {
		c.SSEvent(constant.FAIL, err.Error())
		return
	}
	err = SftpUpload(webFile, "/www-gva/web", "/www-gva/web/deploy.zip", client)
	if err != nil {
		c.SSEvent(constant.FAIL, err.Error())
		return
	}
}

// 环境检测
func (s *ServerApi) Check(c *gin.Context) {

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 将新的上下文传递给需要长时间运行的操作
	var checkRequest request.SSHRequest
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

	Run([]cloud.CmdNode{checkDocker, checkDockerCompose}, client, c)

}

// 安装环境
func (s *ServerApi) Install(c *gin.Context) {

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// 将新的上下文传递给需要长时间运行的操作
	var checkRequest request.SSHRequest
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

	Run([]cloud.CmdNode{docker, dockerCompose}, client, c)

}

func Run(cmds []cloud.CmdNode, client *ssh.Client, c *gin.Context) {
	msg := make(chan cloud.MsgInfo, 100)

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

func SftpUpload(localFile string, dir, remoteFile string, client *ssh.Client) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	// Open the source file
	srcFile, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	//remoteFile所在目录不存在时候 创建目录

	if _, err := sftpClient.Stat(dir); err != nil {
		session, err := client.NewSession()
		if err != nil {
			return err
		}
		err = session.Run("mkdir -p " + dir)
		if err != nil {
			return err
		}
		session.Close()
	}

	// Create the destination file
	dstFile, err := sftpClient.Create(remoteFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the file
	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return err
	}
	return nil
}
