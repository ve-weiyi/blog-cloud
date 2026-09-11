package socialservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteFriendLogic) DeleteFriend(in *socialrpc.DeleteFriendRequest) (*socialrpc.DeleteFriendResponse, error) {
	rows, err := l.svcCtx.TFriendModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	return &socialrpc.DeleteFriendResponse{SuccessCount: rows}, nil
}
