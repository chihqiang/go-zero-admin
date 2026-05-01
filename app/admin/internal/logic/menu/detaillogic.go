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

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.IDRequest) (resp *types.MenuInfo, err error) {
	var menu models.Menu
	if err := l.svcCtx.DB.Where(&models.Menu{ID: req.ID}).First(&menu).Error; err != nil {
		return nil, err
	}

	resp = &types.MenuInfo{}
	if err := copier.Copy(resp, &menu); err != nil {
		return nil, err
	}

	return resp, nil
}