package guestservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/mssola/useragent"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/guestrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/auth/identityx"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/adapter/iplocx"
)

// ipSourceLookup 便于测试替换归属地查询，生产使用百度接口。
var ipSourceLookup = iplocx.GetIpSourceByBaidu

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

	ipSource := ipSourceLookup(ip)
	os := useragent.New(ua).OS()
	browser, _ := useragent.New(ua).Browser()

	// 设备标识由服务端签发：Cookie 或请求头带来的合法值优先，没有才生成新值。
	deviceID, _ := metax.GetDeviceIdFromCtx(l.ctx)
	if identityx.KindOf(deviceID) != identityx.KindGuest {
		deviceID, err = identityx.NewGuestID()
		if err != nil {
			return nil, err
		}
	}

	guest, err := l.svcCtx.TGuestModel.FindOneByDeviceId(l.ctx, deviceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	now := time.Now()
	if guest == nil {
		guest = &model.TGuest{
			DeviceId:  deviceID,
			Os:        os,
			Browser:   browser,
			IpAddress: ip,
			IpSource:  ipSource,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if _, err = l.svcCtx.TGuestModel.Insert(l.ctx, guest); err != nil {
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, err
			}
			// 并发首次访问：另一请求已创建，回查后复用同一行
			guest, err = l.svcCtx.TGuestModel.FindOneByDeviceId(l.ctx, deviceID)
			if err != nil {
				return nil, err
			}
		}
	} else if guest.Os != os || guest.Browser != browser || guest.IpAddress != ip || guest.IpSource != ipSource {
		fields := map[string]interface{}{
			"os":         os,
			"browser":    browser,
			"ip_address": ip,
			"ip_source":  ipSource,
			"updated_at": now,
		}
		if _, err = l.svcCtx.TGuestModel.UpdateFields(l.ctx, fields, "id = ?", guest.Id); err != nil {
			return nil, err
		}

		guest.Os = os
		guest.Browser = browser
		guest.IpAddress = ip
		guest.IpSource = ipSource
		guest.UpdatedAt = now
	}

	return &guestrpc.GetGuestResponse{
		Guest: &guestrpc.Guest{
			Id:        guest.Id,
			DeviceId:  guest.DeviceId,
			Os:        guest.Os,
			Browser:   guest.Browser,
			IpAddress: guest.IpAddress,
			IpSource:  guest.IpSource,
			CreatedAt: guest.CreatedAt.UnixMilli(),
			UpdatedAt: guest.UpdatedAt.UnixMilli(),
		},
	}, nil
}
