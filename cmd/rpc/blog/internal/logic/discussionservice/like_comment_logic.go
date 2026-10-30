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

	comment, err := l.svcCtx.TCommentModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	actorName := l.svcCtx.Notify.ResolveActorName(l.ctx, uid)

	var changed bool
	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		ok, err := likes.Like(l.ctx, likes.Deps{Likes: l.svcCtx.TLikeModel.WithTx(tx)},
			enums.LikeTargetComment, in.Id, uid)
		if err != nil {
			return err
		}
		changed = ok
		if !ok {
			return nil
		}
		// 匿名评论没有收件人，给自己点赞也不通知
		if comment.UserId == "" || comment.UserId == uid {
			return nil
		}

		topicType, topicID := topicOf(comment)
		return l.svcCtx.Notify.Emit(l.ctx, tx, notification.Request{
			Scene:       notify.SceneCommentLiked,
			ActorID:     uid,
			ObjectType:  notify.ObjectComment,
			ObjectID:    in.Id,
			Audience:    notification.AudienceUserIDs,
			AudienceIDs: []string{comment.UserId},
			Payload: notify.Payload{
				ActorName:      actorName,
				CommentExcerpt: excerpt(comment.CommentContent, 60),
				TopicType:      topicType,
				TopicID:        topicID,
			},
		})
	})
	if err != nil {
		return nil, err
	}

	return &discussionrpc.LikeCommentResponse{Liked: true, Changed: changed}, nil
}
