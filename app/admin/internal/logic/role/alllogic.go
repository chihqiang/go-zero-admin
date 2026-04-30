// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"
	"go-zero-admin/app/common/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllLogic {
	return &AllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllLogic) All() (resp []*types.RoleInfo, err error) {
	var roles []*models.Role
	if err := l.svcCtx.DB.Where(&models.Role{}).Order("sort asc").Find(&roles).Error; err != nil {
		return nil, err
	}

	resp = make([]*types.RoleInfo, 0, len(roles))
	for _, role := range roles {
		var info types.RoleInfo
		if err := copier.Copy(&info, role); err != nil {
			return nil, err
		}
		resp = append(resp, &info)
	}

	return resp, nil
}
