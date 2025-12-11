package remark

import (
	"context"

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
		MessageContent: req.MessageContent,
	}

	out, err := l.svcCtx.MessageRpc.AddRemark(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertRemarkTypes(out, nil)
	return resp, nil
}

func ConvertRemarkTypes(req *messagerpc.RemarkDetailsResp, usm map[string]*types.UserInfoVO) (out *types.Remark) {
	return &types.Remark{
		Id:             req.Id,
		UserId:         req.UserId,
		TerminalId:     req.TerminalId,
		MessageContent: req.MessageContent,
		IpAddress:      req.IpAddress,
		IpSource:       req.IpSource,
		IsReview:       req.IsReview,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      req.UpdatedAt,
		User:           usm[req.UserId],
	}
}
