package request

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cloud"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHRequest struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *SSHRequest) Connection(ctx context.Context) (*cloud.Context, error) {
	var cloudCtx cloud.Context
	ctx, cancel := context.WithCancel(ctx)
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	cloudCtx.CTX = ctx
	cloudCtx.Cancel = cancel
	client, err := ssh.Dial("tcp", s.Host+":"+s.Port, config)
	if err != nil {
		return nil, err
	}
	cloudCtx.Client = client
	return &cloudCtx, nil
}
