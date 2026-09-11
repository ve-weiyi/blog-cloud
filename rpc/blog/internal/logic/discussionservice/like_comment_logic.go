package discussionservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type LikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞评论
func (l *LikeCommentLogic) LikeComment(in *discussionrpc.LikeCommentRequest) (*discussionrpc.LikeCommentResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	id := cast.ToString(in.Id)

	entity, err := l.svcCtx.TCommentModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	likeKey := cachekey.GetUserLikeCommentKey(uid)
	countKey := cachekey.CommentLikeCountKey

	ok, err := l.svcCtx.Redis.SIsMember(l.ctx, likeKey, id).Result()
	if err != nil {
		l.Errorf("LikeComment SIsMember error: %v", err)
	}
	if ok {
		entity.LikeCount--
		if err := l.svcCtx.Redis.SRem(l.ctx, likeKey, id).Err(); err != nil {
			l.Errorf("LikeComment SRem error: %v", err)
		}
		if err := l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, -1, id).Err(); err != nil {
			l.Errorf("LikeComment ZIncrBy error: %v", err)
		}
	} else {
		entity.LikeCount++
		if err := l.svcCtx.Redis.SAdd(l.ctx, likeKey, id).Err(); err != nil {
			l.Errorf("LikeComment SAdd error: %v", err)
		}
		if err := l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, 1, id).Err(); err != nil {
			l.Errorf("LikeComment ZIncrBy error: %v", err)
		}
	}

	_, err = l.svcCtx.TCommentModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.LikeCommentResponse{Success: true}, nil
}
