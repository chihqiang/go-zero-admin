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
	var menu models.Menu
	if err := l.svcCtx.DB.Where(&models.Menu{ID: req.ID}).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return err
	}

	var childrenCount int64
	if err := l.svcCtx.DB.Model(&models.Menu{}).Where("pid = ?", req.ID).Count(&childrenCount).Error; err != nil {
		return err
	}
	if childrenCount > 0 {
		return errors.New("存在子菜单，无法删除")
	}

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&menu).Association("Roles").Clear(); err != nil {
			return err
		}
		if err := tx.Delete(&menu).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
