package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesCommentLogic {
	return &DeletesCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除评论
func (l *DeletesCommentLogic) DeletesComment(in *messagerpc.IdsReq) (*messagerpc.BatchResp, error) {
	rows, err := l.svcCtx.TCommentModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &messagerpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
