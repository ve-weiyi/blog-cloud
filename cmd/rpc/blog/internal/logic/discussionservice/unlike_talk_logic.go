package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type UnlikeTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeTalkLogic {
	return &UnlikeTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消点赞说说
func (l *UnlikeTalkLogic) UnlikeTalk(in *discussionrpc.LikeTalkRequest) (*discussionrpc.LikeTalkResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	changed, err := likes.Unlike(l.ctx, l.svcCtx.LikeDeps(), enums.LikeTargetTalk, in.Id, uid)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.LikeTalkResponse{Liked: false, Changed: changed}, nil
}
