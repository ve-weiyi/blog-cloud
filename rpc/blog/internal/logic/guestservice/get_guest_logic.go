package guestservicelogic

import (
	"context"
	"time"

	"github.com/mssola/useragent"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/guestrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/iplocx"
	"github.com/ve-weiyi/vkit/x/cryptox"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
)

type GetGuestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGuestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGuestLogic {
	return &GetGuestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游客环境信息
func (l *GetGuestLogic) GetGuest(in *guestrpc.GetGuestRequest) (*guestrpc.GetGuestResponse, error) {
	ip, err := metax.GetRemoteIPFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	ua, err := metax.GetRemoteAgentFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 分割字符串，提取 IP 部分
	is := iplocx.GetIpSourceByBaidu(ip)
	os := useragent.New(ua).OS()
	browser, _ := useragent.New(ua).Browser()

	did := cryptox.Md5v(ip+os+browser, "")

	// 查找是否已经存在
	vs, _ := l.svcCtx.TGuestModel.FindOneByDeviceId(l.ctx, did)
	if vs == nil {
		vs = &model.TGuest{
			Id:        0,
			DeviceId:  did,
			Os:        os,
			Browser:   browser,
			IpAddress: ip,
			IpSource:  is,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// 插入访客记录
		_, err = l.svcCtx.TGuestModel.Insert(l.ctx, vs)
		if err != nil {
			return nil, err
		}
	}

	return &guestrpc.GetGuestResponse{
		Guest: &guestrpc.Guest{
			Id:        vs.Id,
			DeviceId:  vs.DeviceId,
			Os:        vs.Os,
			Browser:   vs.Browser,
			IpAddress: vs.IpAddress,
			IpSource:  vs.IpSource,
			CreatedAt: vs.CreatedAt.UnixMilli(),
			UpdatedAt: vs.UpdatedAt.UnixMilli(),
		},
	}, nil
}
