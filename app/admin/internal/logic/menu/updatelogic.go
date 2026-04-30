// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

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

func (l *UpdateLogic) Update(req *types.MenuUpdateRequest) (resp *types.MenuInfo, err error) {
	var menu models.Menu
	if err := l.svcCtx.DB.Where(&models.Menu{ID: req.ID}).First(&menu).Error; err != nil {
		return nil, err
	}

	if req.Name != nil {
		menu.Name = *req.Name
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
	if req.Pid != nil {
		menu.Pid = *req.Pid
	}
	if req.Remark != nil {
		menu.Remark = *req.Remark
	}

	if err := l.svcCtx.DB.Save(&menu).Error; err != nil {
		return nil, err
	}

	resp = &types.MenuInfo{}
	if err := copier.Copy(resp, &menu); err != nil {
		return nil, err
	}

	return resp, nil
}
