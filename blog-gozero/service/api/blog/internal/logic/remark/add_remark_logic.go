package remark

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建留言
func NewAddRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRemarkLogic {
	return &AddRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddRemarkLogic) AddRemark(req *types.RemarkNewReq) (resp *types.Remark, err error) {
	in := &messagerpc.RemarkNewReq{
		UserId:         cast.ToString(l.ctx.Value("uid")),
		MessageContent: req.MessageContent,
	}

	out, err := l.svcCtx.MessageRpc.AddRemark(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertRemarkTypes(out)
	return resp, nil
}

func ConvertRemarkTypes(req *messagerpc.RemarkDetails) (out *types.Remark) {
	return &types.Remark{
		Id:             req.Id,
		Nickname:       "",
		Avatar:         "",
		MessageContent: req.MessageContent,
		IpAddress:      req.IpAddress,
		IpSource:       req.IpSource,
		IsReview:       req.IsReview,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      req.UpdatedAt,
	}
}
