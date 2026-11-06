package aiservicelogic

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	roleUser      = "user"
	roleAssistant = "assistant"
)

type AppendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppendMessageLogic {
	return &AppendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 追加消息
func (l *AppendMessageLogic) AppendMessage(in *airpc.AppendMessageRequest) (*airpc.AppendMessageResponse, error) {
	if _, err := findOwnedSession(l.ctx, l.svcCtx, in.SessionId); err != nil {
		return nil, err
	}

	role := strings.TrimSpace(in.Role)
	if role != roleUser && role != roleAssistant {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "消息角色不合法")
	}
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "消息内容不能为空")
	}

	now := time.Now()
	message := &model.TAiMessage{
		SessionId:        in.SessionId,
		Role:             role,
		Content:          content,
		PromptTokens:     maxInt64(in.PromptTokens, 0),
		CompletionTokens: maxInt64(in.CompletionTokens, 0),
		CreatedAt:        now,
	}

	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TAiMessageModel.WithTx(tx).Insert(l.ctx, message); err != nil {
			return err
		}
		_, err := l.svcCtx.TAiSessionModel.WithTx(tx).UpdateFields(l.ctx, map[string]interface{}{
			"message_count":   gorm.Expr("message_count + 1"),
			"last_message_at": now,
			"updated_at":      now,
		}, "id = ?", in.SessionId)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &airpc.AppendMessageResponse{Message: convertMessage(message)}, nil
}

func maxInt64(v, floor int64) int64 {
	if v < floor {
		return floor
	}
	return v
}
