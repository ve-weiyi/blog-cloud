package aiservicelogic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/model"
)

// maxTitleRunes 会话标题的最大字符数。
const maxTitleRunes = 40

func convertSession(in *model.TAiSession) *airpc.Session {
	var lastMessageAt int64
	if in.LastMessageAt != nil {
		lastMessageAt = in.LastMessageAt.UnixMilli()
	}
	return &airpc.Session{
		Id:            in.Id,
		UserId:        in.UserId,
		Title:         in.Title,
		Model:         in.Model,
		MessageCount:  in.MessageCount,
		CreatedAt:     in.CreatedAt.UnixMilli(),
		UpdatedAt:     in.UpdatedAt.UnixMilli(),
		LastMessageAt: lastMessageAt,
	}
}

func convertMessage(in *model.TAiMessage) *airpc.Message {
	return &airpc.Message{
		Id:               in.Id,
		SessionId:        in.SessionId,
		Role:             in.Role,
		Content:          in.Content,
		PromptTokens:     in.PromptTokens,
		CompletionTokens: in.CompletionTokens,
		CreatedAt:        in.CreatedAt.UnixMilli(),
	}
}

// requireUserID 从请求上下文取登录用户；会话接口只对登录用户开放。
func requireUserID(ctx context.Context) (string, error) {
	uid, err := metax.GetUserIdFromCtx(ctx)
	if err != nil || uid == "" {
		return "", bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}
	return uid, nil
}

// findOwnedSession 按 id 查询会话并校验归属，不存在与无权访问统一按不存在处理。
func findOwnedSession(ctx context.Context, svcCtx *svc.ServiceContext, id int64) (*model.TAiSession, error) {
	if id <= 0 {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "会话id不合法")
	}
	uid, err := requireUserID(ctx)
	if err != nil {
		return nil, err
	}

	session, err := svcCtx.TAiSessionModel.FindOne(ctx, "", "id = ? AND user_id = ?", id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "会话不存在")
		}
		return nil, err
	}
	return session, nil
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
