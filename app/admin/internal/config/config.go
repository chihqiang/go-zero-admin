// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"go-zero-admin/pkg/auth"
	"go-zero-admin/pkg/orm"
)

type Config struct {
	rest.RestConf
	Auth auth.JWTConfig
	DB   orm.Config
}
