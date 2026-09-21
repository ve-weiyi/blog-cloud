package discussionservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type LikeTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeTalkLogic {
	return &LikeTalkLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *LikeTalkLogic) LikeTalk(in *discussionrpc.LikeTalkRequest) (*discussionrpc.LikeTalkResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	id := cast.ToString(in.Id)

	entity, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	likeKey := cachekey.GetUserLikeTalkKey(uid)
	countKey := cachekey.TalkLikeCountKey

	ok, err := l.svcCtx.Redis.SIsMember(l.ctx, likeKey, id).Result()
	if err != nil {
		l.Errorf("LikeTalk SIsMember error: %v", err)
	}
	if ok {
		entity.LikeCount--
		if err := l.svcCtx.Redis.SRem(l.ctx, likeKey, id).Err(); err != nil {
			l.Errorf("LikeTalk SRem error: %v", err)
		}
		if err := l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, -1, id).Err(); err != nil {
			l.Errorf("LikeTalk ZIncrBy error: %v", err)
		}
	} else {
		entity.LikeCount++
		if err := l.svcCtx.Redis.SAdd(l.ctx, likeKey, id).Err(); err != nil {
			l.Errorf("LikeTalk SAdd error: %v", err)
		}
		if err := l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, 1, id).Err(); err != nil {
			l.Errorf("LikeTalk ZIncrBy error: %v", err)
		}
	}

	_, err = l.svcCtx.TTalkModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}
	return &discussionrpc.LikeTalkResponse{Success: true}, nil
}
