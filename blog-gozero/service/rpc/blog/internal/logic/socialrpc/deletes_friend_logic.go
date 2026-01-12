package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesFriendLogic {
	return &DeletesFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除友链
func (l *DeletesFriendLogic) DeletesFriend(in *socialrpc.IdsReq) (*socialrpc.BatchResp, error) {
	rows, err := l.svcCtx.TFriendModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &socialrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
