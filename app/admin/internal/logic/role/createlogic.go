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

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.RoleCreateRequest) (resp *types.RoleInfo, err error) {
	var existing models.Role
	if err := l.svcCtx.DB.Where(models.Role{Name: req.Name}).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("角色名称已存在")
	}

	role := models.Role{
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}

	if err := l.svcCtx.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	if len(req.Menus) > 0 {
		var menus []models.Menu
		var menuIDs []int64
		for _, m := range req.Menus {
			menuIDs = append(menuIDs, m.ID)
		}
		l.svcCtx.DB.Where("id IN ?", menuIDs).Find(&menus)
		_ = l.svcCtx.DB.Model(&role).Association("Menus").Replace(&menus)
	}

	resp = &types.RoleInfo{}
	l.svcCtx.DB.Preload("Menus").First(&role, role.ID)
	if err := copier.Copy(resp, &role); err != nil {
		return nil, err
	}

	return resp, nil
}
