package commentrpclogic

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/zero/internal/rpcutil"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *LikeCommentLogic) LikeComment(in *commentrpc.IdReq) (*commentrpc.EmptyResp, error) {
	uid, err := rpcutil.GetRPCUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	id := cast.ToString(in.Id)

	entity, err := l.svcCtx.TCommentModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 设置集合
	likeKey := rediskey.GetUserLikeCommentKey(uid)
	countKey := rediskey.GetCommentLikeCountKey(id)

	ok, _ := l.svcCtx.Redis.SIsMember(l.ctx, likeKey, id).Result()
	if ok {
		// -1
		entity.LikeCount--
		err = l.svcCtx.Redis.SRem(l.ctx, likeKey, id).Err()
		if err != nil {
			return nil, err
		}
		err = l.svcCtx.Redis.Decr(l.ctx, countKey).Err()
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
		err = l.svcCtx.Redis.Incr(l.ctx, countKey).Err()
		if err != nil {
			return nil, err
		}
	}

	_, err = l.svcCtx.TCommentModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}
	return &commentrpc.EmptyResp{}, nil
}
