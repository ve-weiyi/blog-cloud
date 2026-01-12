package comment

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建评论
func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCommentLogic) AddComment(req *types.NewCommentReq) (resp *types.EmptyResp, err error) {
	in := &messagerpc.AddCommentReq{
		ParentId:       req.ParentId,
		TopicId:        req.TopicId,
		ReplyId:        req.ReplyId,
		ReplyUserId:    req.ReplyUserId,
		CommentContent: req.CommentContent,
		Type:           req.Type,
		Status:         req.Status,
	}
	_, err = l.svcCtx.MessageRpc.AddComment(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
