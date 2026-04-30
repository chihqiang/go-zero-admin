// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"errors"
	"go-zero-admin/app/common/models"
	"go-zero-admin/pkg/hash"
	"gorm.io/gorm"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 查找账号
	var account models.Account
	if err := l.svcCtx.DB.Where(models.Account{Email: req.Email}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账户不存在")
		}
		return nil, err
	}

	// 验证密码
	if !hash.Check(req.Password, account.Password) {
		return nil, errors.New("账号或密码错误")
	}
	token, refreshToken, err := l.svcCtx.Jwts.GenerateTokenPair(account.ID)
	if err != nil {
		return nil, err
	}
	return &types.LoginResponse{
		ID:           account.ID,
		Name:         account.Name,
		Email:        account.Email,
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
	}, nil
}
