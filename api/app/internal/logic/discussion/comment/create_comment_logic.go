package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建评论
func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.EmptyResp, err error) {
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)
	did, _ := metax.GetApiDeviceIdFromCtx(l.ctx)

	_, err = l.svcCtx.DiscussionService.CreateComment(l.ctx, &discussionservice.CreateCommentRequest{
		UserId:         uid,
		DeviceId:       did,
		TopicId:        req.TopicId,
		ParentId:       req.ParentId,
		ReplyId:        req.ReplyId,
		ReplyUserId:    req.ReplyUserId,
		CommentContent: req.CommentContent,
		Type:           req.Type,
		Status:         req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
