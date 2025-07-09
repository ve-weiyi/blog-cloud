package remark

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
)

type FindRemarkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取留言列表
func NewFindRemarkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRemarkListLogic {
	return &FindRemarkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindRemarkListLogic) FindRemarkList(req *types.RemarkQuery) (resp *types.PageResp, err error) {
	in := &messagerpc.FindRemarkListReq{
		Paginate: &messagerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}

	out, err := l.svcCtx.MessageRpc.FindRemarkList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
	}

	// 获取用户信息
	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	var list []*types.RemarkBackVO
	for _, v := range out.List {
		m := ConvertRemarkTypes(v, usm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}

func ConvertRemarkTypes(in *messagerpc.RemarkDetailsResp, usm map[string]*types.UserInfoVO) (out *types.RemarkBackVO) {
	out = &types.RemarkBackVO{
		Id:             in.Id,
		UserId:         in.UserId,
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
			out.User = user
		}
	}

	return
}
