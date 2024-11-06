package api

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 同步api列表
func NewSyncApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncApiListLogic {
	return &SyncApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncApiListLogic) SyncApiList(req *types.SyncApiReq) (resp *types.BatchResp, err error) {
	sp, err := parser.Parse(req.ApiFilePath)
	if err != nil {
		return nil, err
	}

	var list []*permissionrpc.ApiNewReq
	for _, g := range sp.Service.Groups {
		var prefix = g.Annotation.Properties["prefix"]
		var group = g.Annotation.Properties["group"]

		var children []*permissionrpc.ApiNewReq
		for _, r := range g.Routes {
			var traceable int64
			if !strings.Contains(strings.ToUpper(r.Method), "GET") {
				if strings.Contains(r.Path, "find") {
					traceable = 0
				} else {
					traceable = 1
				}
			}

			child := &permissionrpc.ApiNewReq{
				Id:        0,
				ParentId:  0,
				Path:      prefix + r.Path,
				Name:      strings.Trim(r.AtDoc.Text, `"`),
				Method:    strings.ToUpper(r.Method),
				Traceable: traceable,
				IsDisable: 0,
				Children:  nil,
			}
			children = append(children, child)
		}
		root := &permissionrpc.ApiNewReq{
			Id:        0,
			ParentId:  0,
			Path:      group,
			Name:      group,
			Method:    "",
			Traceable: 0,
			IsDisable: 0,
			Children:  children,
		}
		list = append(list, root)
	}

	in := &permissionrpc.SyncApiReq{
		Apis: list,
	}

	out, err := l.svcCtx.PermissionRpc.SyncApiList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.RbacHolder.LoadPolicy()
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
