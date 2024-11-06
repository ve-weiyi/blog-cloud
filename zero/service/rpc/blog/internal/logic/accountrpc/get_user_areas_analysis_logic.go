package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAreasAnalysisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAreasAnalysisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAreasAnalysisLogic {
	return &GetUserAreasAnalysisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户分布区域
func (l *GetUserAreasAnalysisLogic) GetUserAreasAnalysis(in *accountrpc.EmptyReq) (*accountrpc.GetUserAreasAnalysisResp, error) {
	// todo: add your logic here and delete this line

	return &accountrpc.GetUserAreasAnalysisResp{}, nil
}
