package cloud

import (
	"context"
	"golang.org/x/crypto/ssh"
)

type SSH interface {
	Connection(ctx context.Context) (*Context, error)
}

type Context struct {
	Client *ssh.Client
	CTX    context.Context
	Cancel context.CancelFunc
}
