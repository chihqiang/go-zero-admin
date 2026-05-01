// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"go-zero-admin/app/admin/internal/config"
	"go-zero-admin/app/common/models"
	"go-zero-admin/pkg/auth"
	"go-zero-admin/pkg/orm"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Jwts   *auth.JWTS
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
		DB:     orm.MustOpen(&c.DB),
		Jwts:   auth.NewJWTS(&c.Auth),
	}
	models.Migrate(s.DB)
	models.MigrateData(s.DB)
	return s
}
