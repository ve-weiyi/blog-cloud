package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *AddOperationLogLogic) AddOperationLog(in *syslogrpc.OperationLogNewReq) (*syslogrpc.OperationLogDetails, error) {
	entity := convertOperationLogIn(in)

	_, err := l.svcCtx.TOperationLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertOperationLogOut(entity), nil
}

func convertOperationLogIn(in *syslogrpc.OperationLogNewReq) (out *model.TOperationLog) {
	out = &model.TOperationLog{
		Id:             0,
		UserId:         in.UserId,
		Nickname:       in.Nickname,
		IpAddress:      in.IpAddress,
		IpSource:       in.IpSource,
		OptModule:      in.OptModule,
		OptDesc:        in.OptDesc,
		RequestUrl:     in.RequestUrl,
		RequestMethod:  in.RequestMethod,
		RequestHeader:  in.RequestHeader,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
		//CreatedAt:      time.Unix(in.CreatedAt, 0),
		//UpdatedAt:      time.Unix(in.UpdatedAt, 0),
	}

	return out
}

func convertOperationLogOut(in *model.TOperationLog) (out *syslogrpc.OperationLogDetails) {
	out = &syslogrpc.OperationLogDetails{
		Id:             in.Id,
		UserId:         in.UserId,
		Nickname:       in.Nickname,
		IpAddress:      in.IpAddress,
		IpSource:       in.IpSource,
		OptModule:      in.OptModule,
		OptDesc:        in.OptDesc,
		RequestUrl:     in.RequestUrl,
		RequestMethod:  in.RequestMethod,
		RequestHeader:  in.RequestHeader,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
		CreatedAt:      in.CreatedAt.Unix(),
		UpdatedAt:      in.UpdatedAt.Unix(),
	}

	return out
}
