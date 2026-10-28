package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
)

type GetCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentLogic {
	return &GetCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询评论
func (l *GetCommentLogic) GetComment(in *discussionrpc.GetCommentRequest) (*discussionrpc.GetCommentResponse, error) {
	entity, err := l.svcCtx.TCommentModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	lm, lcm, err := likeStats(l.ctx, l.svcCtx, enums.LikeTargetComment, []int64{entity.Id})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.GetCommentResponse{Comment: convertCommentOut(entity, lm, lcm)}, nil
}
