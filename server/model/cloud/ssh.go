package cloud

import (
	"golang.org/x/crypto/ssh"
)

type MsgInfo struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

type SSH interface {
	Connection() (*ssh.Client, error)
}

type CmdNode struct {
	Cmd string
}
