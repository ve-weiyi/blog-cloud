package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type LikeTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeTalkLogic {
	return &LikeTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞说说
func (l *LikeTalkLogic) LikeTalk(in *discussionrpc.LikeTalkRequest) (*discussionrpc.LikeTalkResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id); err != nil {
		return nil, err
	}

	err = likes.Toggle(l.ctx, likes.Deps{
		Redis:  l.svcCtx.Redis,
		DB:     l.svcCtx.GormDB,
		Outbox: l.svcCtx.Outbox,
	}, enums.LikeTargetTalk, in.Id, uid,
		cachekey.GetUserLikeTalkKey(uid), cachekey.TalkLikeCountKey)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.LikeTalkResponse{Success: true}, nil
}
