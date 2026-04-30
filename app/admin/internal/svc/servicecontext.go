// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"go-zero-admin/app/admin/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
