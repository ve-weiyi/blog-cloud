package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRemarkLogic {
	return &AddRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建留言
func (l *AddRemarkLogic) AddRemark(in *messagerpc.RemarkNewReq) (*messagerpc.RemarkDetailsResp, error) {
	uid, _ := rpcutils.GetUserIdFromCtx(l.ctx)
	tid, _ := rpcutils.GetTerminalIdFromCtx(l.ctx)
	ip, _ := rpcutils.GetRemoteIPFromCtx(l.ctx)
	is := ipx.GetIpSourceByBaidu(ip)

	entity := &model.TRemark{
		Id:             0,
		UserId:         uid,
		TerminalId:     tid,
		MessageContent: in.MessageContent,
		IpAddress:      ip,
		IpSource:       is,
		IsReview:       l.svcCtx.Config.DefaultRemarkReviewStatus,
		//CreatedAt:      time.Time{},
		//UpdatedAt:      time.Time{},
	}

	_, err := l.svcCtx.TRemarkModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertRemarkOut(entity), nil
}
