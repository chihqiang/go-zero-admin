// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

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

func (l *UpdateLogic) Update(req *types.AccountUpdateRequest) (resp *types.AccountInfo, err error) {
	var account models.Account
	if err := l.svcCtx.DB.Where(&models.Account{ID: req.ID}).Preload("Roles").First(&account).Error; err != nil {
		return nil, err
	}
	_ = copier.Copy(&account, req)
	if err := l.svcCtx.DB.Save(&account).Error; err != nil {
		return nil, err
	}
	resp = &types.AccountInfo{}
	_ = copier.Copy(&resp, &account)
	return resp, nil
}
