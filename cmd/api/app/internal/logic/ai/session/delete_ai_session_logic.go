package session

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/aiservice"
)

type DeleteAiSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除会话
func NewDeleteAiSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAiSessionLogic {
	return &DeleteAiSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAiSessionLogic) DeleteAiSession(req *types.DeleteAiSessionReq) (resp *types.AiSessionDeleteResp, err error) {
	out, err := l.svcCtx.AiService.DeleteSession(l.ctx, &aiservice.DeleteSessionRequest{Id: req.SessionId})
	if err != nil {
		return nil, err
	}

	return &types.AiSessionDeleteResp{SuccessCount: out.GetSuccessCount()}, nil
}
