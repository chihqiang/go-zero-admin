// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"go-zero-admin/app/common/models"
	"go-zero-admin/pkg/hash"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"

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

func (l *CreateLogic) Create(req *types.AccountCreateRequest) (resp *types.AccountInfo, err error) {
	// 检查邮箱
	var existing models.Account
	if err := l.svcCtx.DB.Where(models.Account{Email: req.Email}).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("邮箱已存在")
	}
	// 加密密码
	hashedPassword := hash.Make(req.Password)
	account := models.Account{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   req.Status,
	}
	if err := l.svcCtx.DB.Create(&account).Error; err != nil {
		return nil, err
	}
	// 关联角色
	if len(req.Roles) > 0 {
		var roles []models.Role
		l.svcCtx.DB.Where("id IN ?", req.Roles).Find(&roles)
		_ = l.svcCtx.DB.Model(&account).Association("Roles").Replace(roles)
	}
	resp = &types.AccountInfo{}
	l.svcCtx.DB.Preload("Roles").First(&account, account.ID)
	_ = copier.Copy(&resp, &account)
	return resp, nil
}
