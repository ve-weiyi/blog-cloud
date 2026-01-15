package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOperationLogLogic {
	return &AddOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建操作记录
func (l *AddOperationLogLogic) AddOperationLog(in *syslogrpc.AddOperationLogReq) (*syslogrpc.AddOperationLogResp, error) {
	entity := &model.TOperationLog{
		Id:             0,
		UserId:         in.UserId,
		TerminalId:     in.TerminalId,
		OptModule:      in.OptModule,
		OptDesc:        in.OptDesc,
		RequestUri:     in.RequestUri,
		RequestMethod:  in.RequestMethod,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
	}

	_, err := l.svcCtx.TOperationLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.AddOperationLogResp{
		OperationLog: convertOperationLogOut(entity),
	}, nil
}
