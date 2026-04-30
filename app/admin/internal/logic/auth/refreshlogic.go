// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新令牌
func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshLogic) Refresh(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	token, refreshToken, err := l.svcCtx.Jwts.RefreshTokenPair(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &types.RefreshTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
	}, nil
}
