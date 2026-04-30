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

type AssociateMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssociateMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssociateMenusLogic {
	return &AssociateMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssociateMenusLogic) AssociateMenus(req *types.RoleAssociateMenusRequest) (resp *types.RoleInfo, err error) {
	var role models.Role
	if err := l.svcCtx.DB.Where(&models.Role{ID: req.ID}).Preload("Menus").First(&role).Error; err != nil {
		return nil, err
	}

	var menus []models.Menu
	if len(req.MenuIds) > 0 {
		l.svcCtx.DB.Where("id IN ?", req.MenuIds).Find(&menus)
	}

	if err := l.svcCtx.DB.Model(&role).Association("Menus").Replace(&menus); err != nil {
		return nil, err
	}

	resp = &types.RoleInfo{}
	l.svcCtx.DB.Preload("Menus").First(&role, req.ID)
	if err := copier.Copy(resp, &role); err != nil {
		return nil, err
	}

	return resp, nil
}
