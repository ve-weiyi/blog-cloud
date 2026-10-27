package discussionservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
)

type GetUserLikeTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeTalkLogic {
	return &GetUserLikeTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetUserLikeTalkLogic) GetUserLikeTalk(in *discussionrpc.GetUserLikeTalkRequest) (*discussionrpc.GetUserLikeTalkResponse, error) {
	result, err := l.svcCtx.Redis.SMembers(l.ctx, cachekey.GetUserLikeTalkKey(in.UserId)).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(result))
	for _, v := range result {
		ids = append(ids, cast.ToInt64(v))
	}

	return &discussionrpc.GetUserLikeTalkResponse{LikeTalkIds: ids}, nil
}
