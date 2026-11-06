package aiservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	defaultMessageLimit = 50
	maxMessageLimit     = 200
)

type ListMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessagesLogic {
	return &ListMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询会话消息
func (l *ListMessagesLogic) ListMessages(in *airpc.ListMessagesRequest) (*airpc.ListMessagesResponse, error) {
	if _, err := findOwnedSession(l.ctx, l.svcCtx, in.SessionId); err != nil {
		return nil, err
	}

	limit := int(in.Limit)
	if limit <= 0 {
		limit = defaultMessageLimit
	}
	if limit > maxMessageLimit {
		limit = maxMessageLimit
	}

	records, err := loadSessionMessages(l.ctx, l.svcCtx.TAiMessageModel, in.SessionId, in.After, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*airpc.Message, 0, len(records))
	for _, record := range records {
		list = append(list, convertMessage(record))
	}

	return &airpc.ListMessagesResponse{List: list}, nil
}

// messageLister 是 ListMessages 需要的最小查询能力，便于用假实现验证取数方向。
type messageLister interface {
	FindListAndTotal(ctx context.Context, page int, size int, sorts string, conditions string, args ...interface{}) ([]*model.TAiMessage, int64, error)
}

// loadSessionMessages 取会话消息：after > 0 时增量升序；否则取最近 N 条并翻回正序。
func loadSessionMessages(ctx context.Context, lister messageLister, sessionID, after int64, limit int) ([]*model.TAiMessage, error) {
	if after > 0 {
		records, _, err := lister.FindListAndTotal(ctx, 1, limit, "id asc",
			"session_id = ? AND created_at > ?", sessionID, time.UnixMilli(after))
		return records, err
	}

	records, _, err := lister.FindListAndTotal(ctx, 1, limit, "id desc", "session_id = ?", sessionID)
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}
	return records, nil
}
