// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

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
	// 1. 【查询放在事务外】
	var account models.Account
	if err := l.svcCtx.DB.Where(&models.Account{ID: req.ID}).First(&account).Error; err != nil {
		return err
	}
	// 2. 【获取当前登录用户ID】
	userId, err := l.svcCtx.Jwts.GetID(l.ctx)
	if err != nil {
		return err
	}
	// 3. 不能删除自己
	if account.ID == userId {
		return errors.New("不能删除自己")
	}
	// 4. 【删除操作：事务闭包】
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 删除多对多关联（中间表记录）
		if err := tx.Model(&account).Association("Roles").Clear(); err != nil {
			return err
		}
		// 物理删除用户
		if err := tx.Delete(&account).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
