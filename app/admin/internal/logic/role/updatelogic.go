// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"
	"fmt"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"
	"go-zero-admin/app/common/models"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.RoleUpdateRequest) (resp *types.RoleInfo, err error) {
	var role models.Role
	if err := l.svcCtx.DB.Where(&models.Role{ID: req.ID}).Preload("Menus").First(&role).Error; err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != role.Name {
		var existing models.Role
		if err := l.svcCtx.DB.Where("name = ? AND id != ?", *req.Name, req.ID).First(&existing).Error; err == nil {
			return nil, fmt.Errorf("角色名称已存在")
		}
		role.Name = *req.Name
	}
	if req.Sort != nil {
		role.Sort = *req.Sort
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if req.Remark != nil {
		role.Remark = *req.Remark
	}

	if err := l.svcCtx.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	if req.Menus != nil {
		var menus []models.Menu
		var menuIDs []int64
		for _, m := range req.Menus {
			menuIDs = append(menuIDs, m.ID)
		}
		l.svcCtx.DB.Where("id IN ?", menuIDs).Find(&menus)
		_ = l.svcCtx.DB.Model(&role).Association("Menus").Replace(&menus)
	}

	resp = &types.RoleInfo{}
	l.svcCtx.DB.Preload("Menus").First(&role, req.ID)
	if err := copier.Copy(resp, &role); err != nil {
		return nil, err
	}

	return resp, nil
}