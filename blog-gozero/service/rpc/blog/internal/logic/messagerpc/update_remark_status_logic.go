package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRemarkStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRemarkStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRemarkStatusLogic {
	return &UpdateRemarkStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新留言状态
func (l *UpdateRemarkStatusLogic) UpdateRemarkStatus(in *messagerpc.UpdateRemarkStatusReq) (*messagerpc.BatchResp, error) {
	rows, err := l.svcCtx.TRemarkModel.Updates(l.ctx, map[string]interface{}{
		"status": in.Status,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &messagerpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
