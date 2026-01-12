package socialrpclogic

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeTalkLogic {
	return &LikeTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞说说
func (l *LikeTalkLogic) LikeTalk(in *socialrpc.LikeTalkReq) (*socialrpc.LikeTalkResp, error) {
	uid, err := rpcutils.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	id := cast.ToString(in.Id)

	entity, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	// 设置集合
	likeKey := rediskey.GetUserLikeTalkKey(uid)
	countKey := rediskey.GetTalkLikeCountKey()

	ok, _ := l.svcCtx.Redis.SIsMember(l.ctx, likeKey, id).Result()
	if ok {
		// -1
		entity.LikeCount--
		err = l.svcCtx.Redis.SRem(l.ctx, likeKey, id).Err()
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, -1, id).Err()
		if err != nil {
			return nil, err
		}
	} else {
		// +1
		entity.LikeCount++
		err = l.svcCtx.Redis.SAdd(l.ctx, likeKey, id).Err()
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, 1, id).Err()
		if err != nil {
			return nil, err
		}
	}

	_, err = l.svcCtx.TTalkModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}
	return &socialrpc.LikeTalkResp{}, nil
}
