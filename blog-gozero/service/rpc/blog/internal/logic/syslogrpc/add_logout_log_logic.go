package syslogrpclogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogoutLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddLogoutLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogoutLogLogic {
	return &AddLogoutLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新登录记录
func (l *AddLogoutLogLogic) AddLogoutLog(in *syslogrpc.AddLogoutLogReq) (*syslogrpc.AddLogoutLogResp, error) {
	app, err := rpcutils.GetAppNameFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	exists, _, err := l.svcCtx.TLoginLogModel.FindListAndTotal(l.ctx, 1, 1, "id desc", "user_id = ? and app_name = ?", in.UserId, app)
	if err != nil {
		return nil, err
	}

	if len(exists) == 0 {
		return nil, nil
	}

	exist := exists[0]
	exist.LogoutAt = sql.NullTime{Time: time.Unix(in.LogoutAt, 0), Valid: true}
	exist.UpdatedAt = time.Now()
	_, err = l.svcCtx.TLoginLogModel.Save(l.ctx, exist)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.AddLogoutLogResp{}, nil
}
