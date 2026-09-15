package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type GetMeLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户点赞集合
func NewGetMeLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeLikeLogic {
	return &GetMeLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeLikeLogic) GetMeLike(req *types.GetMeLikeReq) (resp *types.GetMeLikeResp, err error) {
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)

	articleResp, err := l.svcCtx.ContentService.GetUserLikeArticle(l.ctx, &contentservice.GetUserLikeArticleRequest{
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

	talkResp, err := l.svcCtx.DiscussionService.GetUserLikeTalk(l.ctx, &discussionservice.GetUserLikeTalkRequest{
		UserId: uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetMeLikeResp{
		ArticleLikeSet: articleResp.LikeArticleIds,
		CommentLikeSet: commentResp.LikeCommentIds,
		TalkLikeSet:    talkResp.LikeTalkIds,
	}, nil
}
