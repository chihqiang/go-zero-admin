// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-admin/app/admin/internal/svc"
	"go-zero-admin/app/admin/internal/types"
	"go-zero-admin/app/common/models"
	"go-zero-admin/pkg/orm"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.RoleListRequest) (resp *types.RoleListResponse, err error) {
	db := l.svcCtx.DB.Model(&models.Role{}).Order("sort asc")

	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}

	pageData, err := orm.Paginate[*models.Role](db.Preload("Menus"), req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	resp = &types.RoleListResponse{
		PageResponse: types.PageResponse{
			Total: pageData.Total,
			Page:  pageData.Page,
			Size:  pageData.Size,
		},
		Data: make([]*types.RoleInfo, 0, len(pageData.Data)),
	}

	for _, role := range pageData.Data {
		var info types.RoleInfo
		if err := copier.Copy(&info, role); err != nil {
			return nil, err
		}
		resp.Data = append(resp.Data, &info)
	}

	return resp, nil
}
