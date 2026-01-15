package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户点赞列表
func NewGetUserLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeLogic {
	return &GetUserLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLikeLogic) GetUserLike(req *types.EmptyReq) (resp *types.UserLikeResp, err error) {
	uid := cast.ToString(l.ctx.Value(bizheader.HeaderUid))
	articles, err := l.svcCtx.ArticleRpc.FindUserLikeArticle(l.ctx, &articlerpc.FindUserLikeArticleReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	comments, err := l.svcCtx.MessageRpc.FindUserLikeComment(l.ctx, &messagerpc.FindUserLikeCommentReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	talks, err := l.svcCtx.SocialRpc.FindUserLikeTalk(l.ctx, &socialrpc.FindUserLikeTalkReq{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	articleLike := make([]int64, 0)
	commentLike := make([]int64, 0)
	talkLike := make([]int64, 0)
	for _, v := range articles.LikeArticleList {
		articleLike = append(articleLike, v)
	}
	for _, v := range comments.LikeCommentList {
		commentLike = append(commentLike, v)
	}

	for _, v := range talks.LikeTalkList {
		talkLike = append(talkLike, v)
	}

	resp = &types.UserLikeResp{
		ArticleLikeSet: articleLike,
		CommentLikeSet: commentLike,
		TalkLikeSet:    talkLike,
	}

	return
}
