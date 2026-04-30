// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssociateMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssociateMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssociateMenusLogic {
	return &AssociateMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssociateMenusLogic) AssociateMenus(req *types.RoleAssociateMenusRequest) (resp *types.RoleInfo, err error) {
	// todo: add your logic here and delete this line

	return
}
