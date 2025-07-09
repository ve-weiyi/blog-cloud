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

func (l *AddCommentLogic) AddComment(req *types.CommentNewReq) (resp *types.Comment, err error) {
	in := &messagerpc.CommentNewReq{
		ParentId:       req.ParentId,
		TopicId:        req.TopicId,
		ReplyMsgId:     req.ReplyMsgId,
		ReplyUserId:    req.ReplyUserId,
		CommentContent: req.CommentContent,
		Type:           req.Type,
		Status:         req.Status,
	}
	out, err := l.svcCtx.MessageRpc.AddComment(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertCommentTypes(out, nil)
	return resp, nil
}
