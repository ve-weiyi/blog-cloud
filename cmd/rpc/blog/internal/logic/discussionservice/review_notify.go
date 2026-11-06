package discussionservicelogic

import (
	"context"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// emitReviewNotifications 在审核状态迁移时发出通知，必须在审核写入的同一事务里调用。
//
// 通过：给被回复者发「评论被回复」、给内容作者发「内容收到评论」；
// 两人是同一人时只发前者，评论者是作者本人时不通知自己。
// 驳回：先回收此前发出的回复通知（复核驳回后评论对读者不可见），再给作者发审核结果。
func emitReviewNotifications(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, comment *model.TComment, approved bool) error {
	if !approved {
		return emitRejected(ctx, svcCtx, tx, comment)
	}
	return emitApproved(ctx, svcCtx, tx, comment)
}

func emitApproved(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, comment *model.TComment) error {
	topicType, topicID := topicOf(comment)
	actorName := "访客"
	if comment.UserId != "" {
		actorName = svcCtx.Notify.ResolveActorName(ctx, comment.UserId)
	}
	payload := notify.Payload{
		ActorName:      actorName,
		CommentExcerpt: excerpt(comment.CommentContent, 60),
		TopicType:      topicType,
		TopicID:        topicID,
	}

	repliedUID := repliedUserID(ctx, svcCtx, tx, comment)
	if repliedUID != "" && repliedUID != comment.UserId {
		if err := svcCtx.Notify.Emit(ctx, tx, notification.Request{
			Scene:       notify.SceneCommentReplied,
			ActorID:     comment.UserId,
			ObjectType:  notify.ObjectComment,
			ObjectID:    comment.Id,
			Audience:    notification.AudienceUserIDs,
			AudienceIDs: []string{repliedUID},
			Payload:     payload,
		}); err != nil {
			return err
		}
	}

	authorUID, title := contentAuthor(ctx, svcCtx, tx, topicType, topicID)
	if authorUID == "" || authorUID == comment.UserId || authorUID == repliedUID {
		return nil
	}

	payload.ContentTitle = title
	return svcCtx.Notify.Emit(ctx, tx, notification.Request{
		Scene:       notify.SceneContentCommented,
		ActorID:     comment.UserId,
		ObjectType:  topicType,
		ObjectID:    topicID,
		Audience:    notification.AudienceUserIDs,
		AudienceIDs: []string{authorUID},
		Payload:     payload,
	})
}

func emitRejected(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, comment *model.TComment) error {
	// 先回收：复核驳回后评论对读者不可见，之前发出的「回复了你」必须一并消失
	if err := svcCtx.Notify.TerminateObject(ctx, tx, notify.ObjectComment, []int64{comment.Id}); err != nil {
		return err
	}

	// 匿名评论没有收件人；作者本人就是评论者时也不通知
	if comment.UserId == "" {
		return nil
	}

	topicType, topicID := topicOf(comment)
	_, title := contentAuthor(ctx, svcCtx, tx, topicType, topicID)

	// 「评论了你的内容」挂在主题上（按主题聚合），只能按互动定位；
	// 该作者在该主题下还有别的已通过评论时，那条通知仍然成立，不能整组撤掉
	if topicType != "" && !hasOtherApprovedComment(ctx, svcCtx, tx, comment) {
		if err := svcCtx.Notify.RetractInteraction(ctx, tx,
			notify.SceneContentCommented, topicType, topicID, comment.UserId); err != nil {
			return err
		}
	}

	return svcCtx.Notify.Emit(ctx, tx, notification.Request{
		Scene:       notify.SceneContentReviewed,
		ObjectType:  notify.ObjectComment,
		ObjectID:    comment.Id,
		Audience:    notification.AudienceUserIDs,
		AudienceIDs: []string{comment.UserId},
		Payload: notify.Payload{
			Result:       "rejected",
			ContentTitle: title,
			TopicType:    topicType,
			TopicID:      topicID,
		},
	})
}

// hasOtherApprovedComment 判断该作者在同一主题下是否还有别的已通过评论。
//
// 查不出来时按「还有」处理：宁可少撤一条通知，也不要撤掉仍在展示的互动 ——
// 这里不能返回错误，驳回本身比通知清理重要，回滚整笔审核代价更大。
func hasOtherApprovedComment(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, comment *model.TComment) bool {
	count, err := svcCtx.TCommentModel.WithTx(tx).FindCount(ctx,
		"type = ? AND topic_id = ? AND user_id = ? AND is_review = ? AND id <> ?",
		comment.Type, comment.TopicId, comment.UserId, enums.CommentIsReviewYes, comment.Id)
	return err != nil || count > 0
}

// repliedUserID 返回被回复者；没有回复对象时取父评论作者。
func repliedUserID(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, comment *model.TComment) string {
	if comment.ReplyUserId != "" {
		return comment.ReplyUserId
	}
	if comment.ParentId <= 0 {
		return ""
	}
	parent, err := svcCtx.TCommentModel.WithTx(tx).FindById(ctx, comment.ParentId)
	if err != nil {
		return ""
	}
	return parent.UserId
}

// contentAuthor 返回评论所挂内容（文章/说说）的作者与标题；友链评论没有作者。
func contentAuthor(ctx context.Context, svcCtx *svc.ServiceContext, tx *gorm.DB, topicType string, topicID int64) (string, string) {
	switch topicType {
	case notify.ObjectArticle:
		row, err := svcCtx.TArticleModel.WithTx(tx).FindById(ctx, topicID)
		if err != nil {
			return "", ""
		}
		return row.UserId, row.ArticleTitle
	case notify.ObjectTalk:
		row, err := svcCtx.TTalkModel.WithTx(tx).FindById(ctx, topicID)
		if err != nil {
			return "", ""
		}
		return row.UserId, excerpt(row.Content, 40)
	default:
		return "", ""
	}
}
