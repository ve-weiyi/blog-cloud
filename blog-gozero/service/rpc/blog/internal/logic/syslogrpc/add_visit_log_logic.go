package syslogrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type AddVisitLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVisitLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVisitLogLogic {
	return &AddVisitLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建访问记录
func (l *AddVisitLogLogic) AddVisitLog(in *syslogrpc.NewVisitLogReq) (*syslogrpc.EmptyResp, error) {
	uid, _ := rpcutils.GetUserIdFromCtx(l.ctx)
	tid, _ := rpcutils.GetTerminalIdFromCtx(l.ctx)

	entity := &model.TVisitLog{
		Id:         0,
		UserId:     uid,
		TerminalId: tid,
		PageName:   in.PageName,
	}

	_, err := l.svcCtx.TVisitLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.EmptyResp{}, nil
}
