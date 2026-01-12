package tag

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindTagListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取标签列表
func NewFindTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTagListLogic {
	return &FindTagListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindTagListLogic) FindTagList(req *types.QueryTagReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindTagListReq{
		Paginate: &articlerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}
	out, err := l.svcCtx.ArticleRpc.FindTagList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Tag, 0)
	for _, v := range out.List {
		list = append(list, &types.Tag{
			Id:           v.Id,
			TagName:      v.TagName,
			ArticleCount: v.ArticleCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}

	_, err = l.svcCtx.SyslogRpc.AddVisitLog(l.ctx, &syslogrpc.AddVisitLogReq{
		PageName: "标签",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
