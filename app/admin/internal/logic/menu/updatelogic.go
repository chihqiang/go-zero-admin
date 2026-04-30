// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"
	"errors"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"
	"go-zero-admin/app/common/models"

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

func (l *UpdateLogic) Update(req *types.MenuUpdateRequest) (resp *types.MenuInfo, err error) {
	var menu models.Menu
	if err := l.svcCtx.DB.Where(&models.Menu{ID: req.ID}).First(&menu).Error; err != nil {
		return nil, err
	}

	// 检查名称唯一性
	if req.Name != nil && *req.Name != menu.Name {
		var count int64
		l.svcCtx.DB.Model(&models.Menu{}).Where("name = ? AND id != ?", *req.Name, req.ID).Count(&count)
		if count > 0 {
			return nil, errors.New("菜单名称已存在")
		}
		menu.Name = *req.Name
	}

	// 检查父菜单是否存在
	if req.Pid != nil {
		if *req.Pid > 0 {
			var parentMenu models.Menu
			if err := l.svcCtx.DB.First(&parentMenu, *req.Pid).Error; err != nil {
				return nil, errors.New("父菜单不存在")
			}
		}
		menu.Pid = *req.Pid
	}

	if req.MenuType != nil {
		menu.MenuType = *req.MenuType
	}
	if req.Path != nil {
		menu.Path = *req.Path
	}
	if req.Component != nil {
		menu.Component = *req.Component
	}
	if req.Icon != nil {
		menu.Icon = *req.Icon
	}
	if req.Sort != nil {
		menu.Sort = *req.Sort
	}
	if req.ApiUrl != nil {
		menu.ApiUrl = *req.ApiUrl
	}
	if req.ApiMethod != nil {
		menu.ApiMethod = *req.ApiMethod
	}
	if req.Visible != nil {
		menu.Visible = *req.Visible
	}
	if req.Status != nil {
		menu.Status = *req.Status
	}
	if req.Remark != nil {
		menu.Remark = *req.Remark
	}

	if err := l.svcCtx.DB.Save(&menu).Error; err != nil {
		return nil, err
	}

	resp = &types.MenuInfo{
		ID:        menu.ID,
		Pid:       menu.Pid,
		MenuType:  menu.MenuType,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		ApiUrl:    menu.ApiUrl,
		ApiMethod: menu.ApiMethod,
		Visible:   menu.Visible,
		Status:    menu.Status,
		Remark:    menu.Remark,
	}
	return resp, nil
}
