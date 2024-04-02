package request

import (
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHRequest struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *SSHRequest) Connection() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	client, err := ssh.Dial("tcp", s.Host+":"+s.Port, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
