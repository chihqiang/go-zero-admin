// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"go-zero-admin/app/common/models"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户资料
func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile() (resp *types.ProfileResponse, err error) {
	accountID, err := l.svcCtx.Jwts.GetID(l.ctx)
	if err != nil {
		return nil, err
	}
	var account models.Account
	if err := l.svcCtx.DB.Preload("Roles").Preload("Roles.Menus").First(&account, accountID).Error; err != nil {
		return nil, fmt.Errorf("账户不存在")
	}
	resp = &types.ProfileResponse{}
	_ = copier.Copy(&resp, &account)
	// 收集所有菜单并去重
	menuMap := make(map[int64]*models.Menu)
	for _, role := range account.Roles {
		for _, menu := range role.Menus {
			if _, exists := menuMap[menu.ID]; !exists {
				menuMap[menu.ID] = menu
			}
		}
	}
	for _, menu := range menuMap {
		m := &types.ProfileMenuInfo{}
		_ = copier.Copy(&m, &menu)
		resp.Menus = append(resp.Menus, m)
	}
	return resp, nil
}
