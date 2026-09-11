package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
)

type GetUserLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户点赞集合
func NewGetUserLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeLogic {
	return &GetUserLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLikeLogic) GetUserLike(req *types.GetUserLikeReq) (resp *types.GetUserLikeResp, err error) {
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)

	articleResp, err := l.svcCtx.ArticleService.GetUserLikeArticle(l.ctx, &articleservice.GetUserLikeArticleRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	commentResp, err := l.svcCtx.DiscussionService.GetUserLikeComment(l.ctx, &discussionservice.GetUserLikeCommentRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	talkResp, err := l.svcCtx.SocialService.GetUserLikeTalk(l.ctx, &socialservice.GetUserLikeTalkRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetUserLikeResp{
		ArticleLikeSet: articleResp.LikeArticleIds,
		CommentLikeSet: commentResp.LikeCommentIds,
		TalkLikeSet:    talkResp.LikeTalkIds,
	}, nil
}
