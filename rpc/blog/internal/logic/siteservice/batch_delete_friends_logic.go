package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFriendsLogic {
	return &BatchDeleteFriendsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BatchDeleteFriendsLogic) BatchDeleteFriends(in *siterpc.BatchDeleteFriendsRequest) (*siterpc.BatchDeleteFriendsResponse, error) {
	rows, err := l.svcCtx.TFriendModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	return &siterpc.BatchDeleteFriendsResponse{SuccessCount: rows}, nil
}
