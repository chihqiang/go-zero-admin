// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"
	"errors"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"
	"go-zero-admin/app/common/models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.IDRequest) error {
	var role models.Role
	if err := l.svcCtx.DB.Where(&models.Role{ID: req.ID}).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("角色不存在")
		}
		return err
	}

	var accountCount int64
	if err := l.svcCtx.DB.Model(&models.Account{}).Where("EXISTS (SELECT 1 FROM sys_account_roles WHERE account_id = sys_accounts.id AND role_id = ?)", req.ID).Count(&accountCount).Error; err != nil {
		return err
	}
	if accountCount > 0 {
		return errors.New("存在关联的账户，无法删除")
	}

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).Association("Menus").Clear(); err != nil {
			return err
		}
		if err := tx.Delete(&role).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
