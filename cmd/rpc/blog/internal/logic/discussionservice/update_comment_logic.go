package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type UpdateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新评论
func (l *UpdateCommentLogic) UpdateComment(in *discussionrpc.UpdateCommentRequest) (*discussionrpc.UpdateCommentResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil || uid == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}

	if in.Status != enums.CommentStatusNormal && in.Status != enums.CommentStatusEdited {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "评论状态不合法")
	}

	rows, err := l.svcCtx.TCommentModel.UpdateFields(l.ctx, map[string]interface{}{
		"comment_content": in.CommentContent,
		"status":          in.Status,
	}, "id = ? AND user_id = ?", in.Id, uid)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, bizerr.NewBizError(bizcode.CodeNoPermission, "无权修改该评论")
	}

	return &discussionrpc.UpdateCommentResponse{Success: true}, nil
}
