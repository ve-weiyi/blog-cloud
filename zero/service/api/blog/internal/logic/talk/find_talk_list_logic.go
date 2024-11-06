package talk

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/talkrpc"

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

func (l *FindTalkListLogic) FindTalkList(req *types.TalkQueryReq) (resp *types.PageResp, err error) {
	in := &talkrpc.FindTalkListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	out, err := l.svcCtx.TalkRpc.FindTalkList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var tids []int64
	var uids []string
	for _, v := range out.List {
		tids = append(tids, v.Id)
		uids = append(uids, v.UserId)
	}

	// 查询用户信息
	users, err := l.svcCtx.AccountRpc.FindUserList(l.ctx, &accountrpc.FindUserListReq{
		UserIds: uids,
	})
	if err != nil {
		return nil, err
	}

	usm := make(map[string]*accountrpc.User)
	for _, v := range users.List {
		usm[v.UserId] = v
	}

	// 查询评论量
	counts, err := l.svcCtx.CommentRpc.FindTopicCommentCounts(l.ctx, &commentrpc.IdsReq{
		Ids: tids,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.Talk, 0)
	for _, v := range out.List {
		m := ConvertTalkTypes(v, usm, counts.TopicCommentCounts)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertTalkTypes(in *talkrpc.TalkDetails, usm map[string]*accountrpc.User, csm map[int64]int64) (out *types.Talk) {
	out = &types.Talk{
		Id:           in.Id,
		UserId:       in.UserId,
		Content:      in.Content,
		ImgList:      in.ImgList,
		IsTop:        in.IsTop,
		Status:       in.Status,
		LikeCount:    in.LikeCount,
		CommentCount: 0,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}

	// 用户信息
	if out.UserId != "" {
		user, ok := usm[out.UserId]
		if ok && user != nil {
			out.Nickname = user.Nickname
			out.Avatar = user.Avatar
		}
	}

	// 评论量
	count, ok := csm[out.Id]
	if ok {
		out.CommentCount = count
	}
	return
}
