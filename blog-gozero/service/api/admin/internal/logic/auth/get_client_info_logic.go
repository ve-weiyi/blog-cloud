package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetClientInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取客户端信息
func NewGetClientInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClientInfoLogic {
	return &GetClientInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClientInfoLogic) GetClientInfo(req *types.GetClientInfoReq) (resp *types.GetClientInfoResp, err error) {
	out, err := l.svcCtx.AccountRpc.GetClientInfo(l.ctx, &accountrpc.GetClientInfoReq{})
	if err != nil {
		return nil, err
	}

	return &types.GetClientInfoResp{
		Id:         out.Visitor.Id,
		TerminalId: out.Visitor.TerminalId,
		Os:         out.Visitor.Os,
		Browser:    out.Visitor.Browser,
		IpAddress:  out.Visitor.IpAddress,
		IpSource:   out.Visitor.IpSource,
	}, nil
}
