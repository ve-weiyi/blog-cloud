package comment

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新评论
func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCommentLogic) UpdateComment(req *types.UpdateCommentReq) (resp *types.Comment, err error) {
	in := &messagerpc.UpdateCommentReq{
		Id:             req.Id,
		ReplyUserId:    req.ReplyUserId,
		CommentContent: req.CommentContent,
		Status:         req.Status,
	}

	out, err := l.svcCtx.MessageRpc.UpdateComment(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertCommentTypes(out, nil)
	return resp, nil
}
