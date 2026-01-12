package messagerpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新评论
func (l *UpdateCommentLogic) UpdateComment(in *messagerpc.UpdateCommentReq) (*messagerpc.CommentDetailsResp, error) {
	uid, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 查找评论
	comment, err := l.svcCtx.TCommentModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if comment.UserId != uid {
		return nil, fmt.Errorf("无权限操作")
	}

	// 更新评论
	comment.CommentContent = in.CommentContent
	comment.Status = in.Status
	_, err = l.svcCtx.TCommentModel.Save(l.ctx, comment)
	if err != nil {
		return nil, err
	}

	return &messagerpc.CommentDetailsResp{}, nil
}
