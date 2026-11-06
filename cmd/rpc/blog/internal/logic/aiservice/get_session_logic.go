package aiservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type GetSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionLogic {
	return &GetSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询会话详情
func (l *GetSessionLogic) GetSession(in *airpc.GetSessionRequest) (*airpc.GetSessionResponse, error) {
	session, err := findOwnedSession(l.ctx, l.svcCtx, in.Id)
	if err != nil {
		return nil, err
	}

	return &airpc.GetSessionResponse{Session: convertSession(session)}, nil
}
