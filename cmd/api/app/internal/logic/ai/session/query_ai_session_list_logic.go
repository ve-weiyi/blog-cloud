package session

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/aiservice"
)

type QueryAiSessionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页查询我的会话
func NewQueryAiSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAiSessionListLogic {
	return &QueryAiSessionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAiSessionListLogic) QueryAiSessionList(req *types.QueryAiSessionListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.AiService.ListSessions(l.ctx, &aiservice.ListSessionsRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.AiSessionVO, 0, len(out.GetList()))
	for _, item := range out.GetList() {
		list = append(list, &types.AiSessionVO{
			Id:            item.GetId(),
			Title:         item.GetTitle(),
			Model:         item.GetModel(),
			MessageCount:  item.GetMessageCount(),
			CreatedAt:     item.GetCreatedAt(),
			UpdatedAt:     item.GetUpdatedAt(),
			LastMessageAt: item.GetLastMessageAt(),
		})
	}

	return &types.ListResult{
		Page:     out.GetListResult().GetPage(),
		PageSize: out.GetListResult().GetPageSize(),
		Total:    out.GetListResult().GetTotal(),
		List:     list,
	}, nil
}
