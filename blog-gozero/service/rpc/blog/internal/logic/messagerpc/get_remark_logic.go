package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRemarkLogic {
	return &GetRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询留言
func (l *GetRemarkLogic) GetRemark(in *messagerpc.IdReq) (*messagerpc.RemarkDetailsResp, error) {
	entity, err := l.svcCtx.TRemarkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return convertRemarkOut(entity), nil
}
