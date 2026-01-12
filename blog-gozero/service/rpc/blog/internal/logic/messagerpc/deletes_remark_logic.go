package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesRemarkLogic {
	return &DeletesRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除留言
func (l *DeletesRemarkLogic) DeletesRemark(in *messagerpc.DeletesRemarkReq) (*messagerpc.DeletesRemarkResp, error) {
	rows, err := l.svcCtx.TRemarkModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &messagerpc.DeletesRemarkResp{
		SuccessCount: rows,
	}, nil
}
