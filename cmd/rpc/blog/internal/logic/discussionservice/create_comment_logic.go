package discussionservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建评论
//
// 被回复者由服务端按被回复评论反查，不采信请求体里的 reply_user_id：
// 这个值既用于通知收件人，也用于「回复 @某人」的展示。
func (l *CreateCommentLogic) CreateComment(in *discussionrpc.CreateCommentRequest) (*discussionrpc.CreateCommentResponse, error) {
	replyUserID, err := resolveReplyUser(l.ctx, l.svcCtx.TCommentModel, in.ReplyId, in.ParentId)
	if errors.Is(err, errReplyTargetNotFound) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "被回复的评论不存在")
	}
	if err != nil {
		return nil, err
	}

	entity := &model.TComment{
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		TopicId:        in.TopicId,
		ParentId:       in.ParentId,
		ReplyId:        in.ReplyId,
		ReplyUserId:    replyUserID,
		CommentContent: in.CommentContent,
		Type:           in.Type,
		Status:         enums.CommentStatusNormal,
		IsReview:       enums.CommentIsReviewNo,
	}

	if _, err := l.svcCtx.TCommentModel.Insert(l.ctx, entity); err != nil {
		return nil, err
	}

	return &discussionrpc.CreateCommentResponse{Id: entity.Id}, nil
}
