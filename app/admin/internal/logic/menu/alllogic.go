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

func (l *AllLogic) All() (resp []*types.MenuAllInfo, err error) {
	var menus []*models.Menu
	if err := l.svcCtx.DB.Where(&models.Menu{}).Find(&menus).Error; err != nil {
		return nil, err
	}
	resp = make([]*types.MenuAllInfo, 0, len(menus))
	for _, menu := range menus {
		var info *types.MenuAllInfo
		copier.Copy(&info, &menu)
		resp = append(resp, info)
	}
	return resp, nil
}
