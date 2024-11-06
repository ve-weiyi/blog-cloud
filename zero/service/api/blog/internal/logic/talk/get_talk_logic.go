package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/talkrpc"
)

type GetTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询说说
func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTalkLogic) GetTalk(req *types.IdReq) (resp *types.Talk, err error) {
	in := &talkrpc.IdReq{
		Id: req.Id,
	}

	out, err := l.svcCtx.TalkRpc.GetTalk(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	users, err := l.svcCtx.AccountRpc.FindUserList(l.ctx, &accountrpc.FindUserListReq{
		UserIds: []string{out.UserId},
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
		Ids: []int64{out.Id},
	})
	if err != nil {
		return nil, err
	}

	resp = ConvertTalkTypes(out, usm, counts.TopicCommentCounts)
	return resp, nil
}
