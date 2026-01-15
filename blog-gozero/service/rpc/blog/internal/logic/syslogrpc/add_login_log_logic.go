package syslogrpclogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type AddLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLoginLogLogic {
	return &AddLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建登录记录
func (l *AddLoginLogLogic) AddLoginLog(in *syslogrpc.AddLoginLogReq) (*syslogrpc.AddLoginLogResp, error) {
	app, _ := rpcutils.GetAppNameFromCtx(l.ctx)
	tid, _ := rpcutils.GetTerminalIdFromCtx(l.ctx)

	now := time.Now()
	entity := &model.TLoginLog{
		Id:         0,
		UserId:     in.UserId,
		TerminalId: tid,
		LoginType:  in.LoginType,
		AppName:    app,
		LoginAt:    now,
		LogoutAt:   sql.NullTime{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err := l.svcCtx.TLoginLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.AddLoginLogResp{
		LoginLog: convertLoginLogOut(entity),
	}, nil
}
