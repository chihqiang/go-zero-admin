// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

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

func (l *ListLogic) List(req *types.MenuListRequest) (resp *types.MenuListResponse, err error) {
	db := l.svcCtx.DB.Model(&models.Menu{}).Order("pid asc, sort asc")

	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status {
		db = db.Where("status = ?", req.Status)
	}

	pageData, err := orm.Paginate[*models.Menu](db, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	resp = &types.MenuListResponse{
		PageResponse: types.PageResponse{
			Total: pageData.Total,
			Page:  pageData.Page,
			Size:  pageData.Size,
		},
		Data: make([]*types.MenuInfo, 0, len(pageData.List)),
	}

	for _, menu := range pageData.List {
		var info types.MenuInfo
		if err := copier.Copy(&info, menu); err != nil {
			return nil, err
		}
		resp.Data = append(resp.Data, &info)
	}

	return resp, nil
}
