package accountrpclogic

import (
	"context"
	"time"

	"github.com/mssola/useragent"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/ipx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetClientInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetClientInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClientInfoLogic {
	return &GetClientInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取客户端信息
func (l *GetClientInfoLogic) GetClientInfo(in *accountrpc.GetClientInfoReq) (*accountrpc.GetClientInfoResp, error) {
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

	terminalId := cryptox.Md5v(ip+os+browser, "")

	// 查找是否已经存在
	vs, _ := l.svcCtx.TVisitorModel.FindOneByTerminalId(l.ctx, terminalId)
	if vs == nil {
		vs = &model.TVisitor{
			Id:         0,
			TerminalId: terminalId,
			Os:         os,
			Browser:    browser,
			IpAddress:  ip,
			IpSource:   is,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		// 插入访客记录
		_, err = l.svcCtx.TVisitorModel.Insert(l.ctx, vs)
		if err != nil {
			return nil, err
		}
	}

	return &accountrpc.GetClientInfoResp{
		Visitor: &accountrpc.VisitorInfo{
			Id:         vs.Id,
			TerminalId: vs.TerminalId,
			Os:         vs.Os,
			Browser:    vs.Browser,
			IpAddress:  vs.IpAddress,
			IpSource:   vs.IpSource,
			CreatedAt:  vs.CreatedAt.UnixMilli(),
			UpdatedAt:  vs.UpdatedAt.UnixMilli(),
		},
	}, nil
}
