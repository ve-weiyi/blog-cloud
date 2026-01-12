package talk

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"

	"github.com/zeromicro/go-zero/core/logx"
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
	}
	out, err := l.svcCtx.SocialRpc.FindTalkList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var tids []int64
	for _, v := range out.List {
		tids = append(tids, v.Id)
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

	// 查询评论量
	counts, err := l.svcCtx.MessageRpc.FindCommentReplyCounts(l.ctx, &messagerpc.FindCommentReplyCountsReq{
		TopicIds: tids,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.Talk, 0)
	for _, v := range out.List {
		m := convertTalkTypes(v, usm, counts.TopicCommentCounts)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}

func convertTalkTypes(in *socialrpc.Talk, usm map[string]*types.UserInfoVO, csm map[int64]int64) (out *types.Talk) {
	out = &types.Talk{
		Id:           in.Id,
		UserId:       in.UserId,
		Content:      in.Content,
		ImgList:      in.ImgList,
		IsTop:        in.IsTop,
		Status:       in.Status,
		LikeCount:    in.LikeCount,
		CommentCount: csm[in.Id],
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		UserInfo:     usm[in.UserId],
	}
	return
}
