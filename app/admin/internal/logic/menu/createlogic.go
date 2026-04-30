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
	var menu models.Menu
	copier.Copy(&menu, req)
	if err := l.svcCtx.DB.Create(&menu).Error; err != nil {
		return nil, err
	}

	resp = &types.MenuInfo{}
	copier.Copy(resp, &menu)
	return resp, nil
}
