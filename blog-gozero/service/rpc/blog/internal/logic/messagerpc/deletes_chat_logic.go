package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesChatLogic {
	return &DeletesChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除聊天记录
func (l *DeletesChatLogic) DeletesChat(in *messagerpc.DeletesChatReq) (*messagerpc.DeletesChatResp, error) {
	rows, err := l.svcCtx.TChatModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &messagerpc.DeletesChatResp{
		SuccessCount: rows,
	}, nil
}
