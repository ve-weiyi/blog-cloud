package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
)

type GetTouristInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取游客身份信息
func NewGetTouristInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTouristInfoLogic {
	return &GetTouristInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTouristInfoLogic) GetTouristInfo(req *types.EmptyReq) (resp *types.GetTouristInfoResp, err error) {
	tourist, err := l.svcCtx.AccountRpc.GetTouristInfo(l.ctx, &accountrpc.GetTouristInfoReq{})
	if err != nil {
		return nil, err
	}

	return &types.GetTouristInfoResp{
		TouristId: tourist.TouristId,
	}, nil
}
