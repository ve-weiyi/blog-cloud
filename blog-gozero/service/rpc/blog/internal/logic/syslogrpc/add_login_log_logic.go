package syslogrpclogic

import (
	"context"
	"database/sql"
	"time"

	"github.com/mssola/useragent"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *AddLoginLogLogic) AddLoginLog(in *syslogrpc.LoginLogNewReq) (*syslogrpc.EmptyResp, error) {
	app, err := rpcutils.GetAppNameFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	ip, err := rpcutils.GetRemoteIPFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	ua, err := rpcutils.GetRemoteAgentFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 分割字符串，提取 IP 部分
	is := ipx.GetIpSourceByBaidu(ip)
	os := useragent.New(ua).OS()
	browser, _ := useragent.New(ua).Browser()

	now := time.Now()
	entity := &model.TLoginLog{
		Id:        0,
		UserId:    in.UserId,
		LoginType: in.LoginType,
		AppName:   app,
		Os:        os,
		Browser:   browser,
		IpAddress: ip,
		IpSource:  is,
		LoginAt:   now,
		LogoutAt:  sql.NullTime{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = l.svcCtx.TLoginLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.EmptyResp{}, nil
}
