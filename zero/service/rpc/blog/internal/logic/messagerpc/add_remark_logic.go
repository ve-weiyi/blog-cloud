package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/internal/rpcutil"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *AddRemarkLogic) AddRemark(in *messagerpc.RemarkNewReq) (*messagerpc.RemarkDetails, error) {
	ip, _ := rpcutil.GetRPCClientIP(l.ctx)
	is, _ := ipx.GetIpSourceByBaidu(ip)

	entity := &model.TRemark{
		Id:             0,
		UserId:         in.UserId,
		MessageContent: in.MessageContent,
		IpAddress:      ip,
		IpSource:       is,
		IsReview:       1,
		//CreatedAt:      time.Time{},
		//UpdatedAt:      time.Time{},
	}

	_, err := l.svcCtx.TRemarkModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertRemarkOut(entity), nil
}

func convertRemarkOut(in *model.TRemark) (out *messagerpc.RemarkDetails) {
	out = &messagerpc.RemarkDetails{
		Id:             in.Id,
		UserId:         in.UserId,
		MessageContent: in.MessageContent,
		IpAddress:      in.IpAddress,
		IpSource:       in.IpSource,
		IsReview:       in.IsReview,
		CreatedAt:      in.CreatedAt.Unix(),
		UpdatedAt:      in.UpdatedAt.Unix(),
	}

	return out
}
