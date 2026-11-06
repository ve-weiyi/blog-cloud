package session

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/aiservice"
)

type QueryAiMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询会话消息
func NewQueryAiMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAiMessageListLogic {
	return &QueryAiMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAiMessageListLogic) QueryAiMessageList(req *types.QueryAiMessageListReq) (resp *types.AiMessageListResp, err error) {
	out, err := l.svcCtx.AiService.ListMessages(l.ctx, &aiservice.ListMessagesRequest{
		SessionId: req.SessionId,
		After:     req.After,
		Limit:     req.Limit,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.AiMessageVO, 0, len(out.GetList()))
	for _, item := range out.GetList() {
		list = append(list, &types.AiMessageVO{
			Id:        item.GetId(),
			SessionId: item.GetSessionId(),
			Role:      item.GetRole(),
			Content:   item.GetContent(),
			CreatedAt: item.GetCreatedAt(),
		})
	}

	return &types.AiMessageListResp{List: list}, nil
}
