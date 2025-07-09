package syslogrpclogic

import (
	"context"

	"github.com/mssola/useragent"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *AddVisitLogLogic) AddVisitLog(in *syslogrpc.VisitLogNewReq) (*syslogrpc.EmptyResp, error) {
	uid, _ := rpcutils.GetUserIdFromCtx(l.ctx)
	tid, _ := rpcutils.GetTerminalIdFromCtx(l.ctx)
	ip, _ := rpcutils.GetRemoteIPFromCtx(l.ctx)
	ua, _ := rpcutils.GetRemoteAgentFromCtx(l.ctx)

	// 分割字符串，提取 IP 部分
	is := ipx.GetIpSourceByBaidu(ip)
	os := useragent.New(ua).OS()
	browser, _ := useragent.New(ua).Browser()

	entity := &model.TVisitLog{
		Id:         0,
		UserId:     uid,
		TerminalId: tid,
		PageName:   in.PageName,
		IpAddress:  ip,
		IpSource:   is,
		Os:         os,
		Browser:    browser,
	}

	_, err := l.svcCtx.TVisitLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.EmptyResp{}, nil
}
