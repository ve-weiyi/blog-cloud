package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnalysisMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnalysisMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalysisMessageLogic {
	return &AnalysisMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 消息数据分析
func (l *AnalysisMessageLogic) AnalysisMessage(in *messagerpc.AnalysisMessageReq) (*messagerpc.AnalysisMessageResp, error) {
	rc, err := l.svcCtx.TRemarkModel.FindCount(l.ctx, "")
	if err != nil {
		return nil, err
	}

	return &messagerpc.AnalysisMessageResp{
		RemarkCount: rc,
	}, nil
}
