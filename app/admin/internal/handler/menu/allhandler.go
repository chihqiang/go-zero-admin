// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"

	xhttp "github.com/zeromicro/x/http"
	"go-zero-admin/app/admin/internal/logic/menu"
	"go-zero-admin/app/admin/internal/svc"
)

func AllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := menu.NewAllLogic(r.Context(), svcCtx)
		resp, err := l.All()
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
