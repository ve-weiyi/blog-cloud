package remark

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新留言
func NewUpdateRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRemarkLogic {
	return &UpdateRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRemarkLogic) UpdateRemark(req *types.RemarkNewReq) (resp *types.RemarkBackDTO, err error) {
	in := &messagerpc.RemarkUpdateReq{
		Id:       req.Id,
		IsReview: req.IsReview,
	}
	out, err := l.svcCtx.MessageRpc.UpdateRemark(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertRemarkTypes(out, nil)
	return resp, nil
}

func ConvertRemarkTypes(in *messagerpc.RemarkDetails, usm map[string]*accountrpc.User) (out *types.RemarkBackDTO) {
	out = &types.RemarkBackDTO{
		Id:             in.Id,
		Nickname:       "",
		Avatar:         "",
		MessageContent: in.MessageContent,
		IpAddress:      in.IpAddress,
		IpSource:       in.IpSource,
		Time:           0,
		IsReview:       in.IsReview,
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
	}

	// 用户信息
	if in.UserId != "" {
		user, ok := usm[in.UserId]
		if ok && user != nil {
			out.Nickname = user.Nickname
			out.Avatar = user.Avatar
		}
	}

	return
}
