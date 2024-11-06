package permissionrpclogic

import (
	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
)

func convertApiIn(in *permissionrpc.ApiNewReq) (out *model.TApi) {
	out = &model.TApi{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Method:    in.Method,
		Traceable: in.Traceable,
		IsDisable: in.IsDisable,
	}

	return out
}

func convertApiOut(in *model.TApi) (out *permissionrpc.ApiDetails) {
	out = &permissionrpc.ApiDetails{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Method:    in.Method,
		Traceable: in.Traceable,
		IsDisable: in.IsDisable,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
	}

	return out
}
