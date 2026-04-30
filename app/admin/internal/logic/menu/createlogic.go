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

func (l *CreateLogic) Create(req *types.MenuCreateRequest) (resp *types.MenuInfo, err error) {
	// 检查名称唯一性
	var count int64
	l.svcCtx.DB.Model(&models.Menu{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("菜单名称已存在")
	}

	// 检查父菜单是否存在
	if req.Pid > 0 {
		var parentMenu models.Menu
		if err := l.svcCtx.DB.First(&parentMenu, req.Pid).Error; err != nil {
			return nil, errors.New("父菜单不存在")
		}
	}

	menu := models.Menu{
		Pid:       req.Pid,
		MenuType:  req.MenuType,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		ApiUrl:    req.ApiUrl,
		ApiMethod: req.ApiMethod,
		Visible:   req.Visible,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	if err := l.svcCtx.DB.Create(&menu).Error; err != nil {
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
