package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
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
func (l *LikeTalkLogic) LikeTalk(in *discussionrpc.LikeTalkRequest) (*discussionrpc.LikeTalkResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	talk, err := l.svcCtx.TTalkModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	actorName := l.svcCtx.Notify.ResolveActorName(l.ctx, uid)

	var changed bool
	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		ok, err := likes.Like(l.ctx, likes.Deps{Likes: l.svcCtx.TLikeModel.WithTx(tx)},
			enums.LikeTargetTalk, in.Id, uid)
		if err != nil {
			return err
		}
		changed = ok
		if !ok {
			return nil
		}
		if talk.UserId == "" || talk.UserId == uid {
			return nil
		}

		return l.svcCtx.Notify.Emit(l.ctx, tx, notification.Request{
			Scene:       notify.SceneTalkLiked,
			ActorID:     uid,
			ObjectType:  notify.ObjectTalk,
			ObjectID:    in.Id,
			Audience:    notification.AudienceUserIDs,
			AudienceIDs: []string{talk.UserId},
			Payload: notify.Payload{
				ActorName: actorName,
				TopicType: notify.ObjectTalk,
				TopicID:   in.Id,
			},
		})
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.LikeTalkResponse{Liked: true, Changed: changed}, nil
}
