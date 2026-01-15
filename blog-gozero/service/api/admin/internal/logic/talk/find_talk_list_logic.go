package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"
)

type FindTalkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取说说列表
func NewFindTalkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTalkListLogic {
	return &FindTalkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindTalkListLogic) FindTalkList(req *types.QueryTalkReq) (resp *types.PageResp, err error) {
	in := &socialrpc.FindTalkListReq{
		Paginate: &socialrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		Status: req.Status,
	}

	out, err := l.svcCtx.SocialRpc.FindTalkList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	usm, err := apiutils.BatchQuery(out.List,
		func(v *socialrpc.Talk) string {
			return v.UserId
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	var list []*types.TalkBackVO
	for _, v := range out.List {
		m := &types.TalkBackVO{
			Id:           v.Id,
			UserId:       v.UserId,
			Content:      v.Content,
			ImgList:      v.ImgList,
			IsTop:        v.IsTop,
			Status:       v.Status,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			UserInfo:     usm[v.UserId],
		}

		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
