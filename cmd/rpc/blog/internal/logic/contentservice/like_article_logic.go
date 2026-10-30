package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type LikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞文章
//
// 点赞与通知事件同事务：事件写失败，点赞一并回滚。取消点赞不发通知，
// 重复点赞由点赞明细的唯一键挡住，影响行数为 0 时不产生事件。
func (l *LikeArticleLogic) LikeArticle(in *contentrpc.LikeArticleRequest) (*contentrpc.LikeArticleResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	article, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	actorName := l.svcCtx.Notify.ResolveActorName(l.ctx, uid)

	var changed bool
	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		ok, err := likes.Like(l.ctx, likes.Deps{Likes: l.svcCtx.TLikeModel.WithTx(tx)},
			enums.LikeTargetArticle, in.Id, uid)
		if err != nil {
			return err
		}
		changed = ok
		if !ok {
			return nil
		}
		// 给自己点赞不通知，作者为空（数据异常）也不通知
		if article.UserId == "" || article.UserId == uid {
			return nil
		}

		return l.svcCtx.Notify.Emit(l.ctx, tx, notification.Request{
			Scene:       notify.SceneArticleLiked,
			ActorID:     uid,
			ObjectType:  notify.ObjectArticle,
			ObjectID:    in.Id,
			Audience:    notification.AudienceUserIDs,
			AudienceIDs: []string{article.UserId},
			Payload: notify.Payload{
				ActorName:    actorName,
				ContentTitle: article.ArticleTitle,
				TopicType:    notify.ObjectArticle,
				TopicID:      in.Id,
			},
		})
	})
	if err != nil {
		return nil, err
	}

	return &contentrpc.LikeArticleResponse{Liked: true, Changed: changed}, nil
}
