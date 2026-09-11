package notificationservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"
)

type ListUserInboxRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserInboxRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserInboxRecordsLogic {
	return &ListUserInboxRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserInboxRecordsLogic) ListUserInboxRecords(in *notificationrpc2.ListUserInboxRecordsRequest) (*notificationrpc2.ListUserInboxRecordsResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}

	opts = append(opts, queryx.WithCondition("channel = ?", "inbox"))
	opts = append(opts, queryx.WithCondition("recipient = ?", in.UserId))
	opts = append(opts, queryx.WithCondition("status != ?", "revoked"))

	if in.OnlyUnread != nil && *in.OnlyUnread == 1 {
		opts = append(opts, queryx.WithCondition("status = ?", "unread"))
	}

	var titleMsgIds []int64
	if in.Title != nil && *in.Title != "" {
		titleMsgIds = l.getMessageIdsByTitle(*in.Title)
		if len(titleMsgIds) == 0 {
			return &notificationrpc.ListUserInboxRecordsResponse{
				PageResult:  &notificationrpc.PageResult{Page: 1, PageSize: 10, Total: 0},
				UnreadTotal: 0,
			}, nil
		}
		placeholders := make([]string, len(titleMsgIds))
		args := make([]any, len(titleMsgIds))
		for i, id := range titleMsgIds {
			placeholders[i] = "?"
			args[i] = id
		}
		opts = append(opts, queryx.WithCondition("message_id IN ("+strings.Join(placeholders, ",")+")", args...))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TNotifyRecordModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	// 收集 message_id 并批量查询消息
	msgMap := l.getMessageMap(records)

	var list []*notificationrpc.NotifyRecord
	for _, v := range records {
		d := convertTNotifyRecordToProto(v)
		if msg, ok := msgMap[v.MessageId]; ok {
			d.Title = msg.Title
			d.Category = msg.Category
			d.Level = msg.Level
		}
		list = append(list, d)
	}

	// 查询未读数量
	var unreadTotal int64
	if in.OnlyUnread != nil && *in.OnlyUnread == 1 {
		// 只查未读时，total 即未读数
		unreadTotal = total
	} else if len(titleMsgIds) > 0 {
		placeholders := make([]string, len(titleMsgIds))
		args := make([]any, 0, 3+len(titleMsgIds))
		args = append(args, "inbox", in.UserId, "unread")
		for i, id := range titleMsgIds {
			placeholders[i] = "?"
			args = append(args, id)
		}
		unreadTotal, err = l.svcCtx.TNotifyRecordModel.FindCount(l.ctx,
			"channel = ? AND recipient = ? AND status = ? AND message_id IN ("+strings.Join(placeholders, ",")+")", args...)
		if err != nil {
			l.Errorf("ListUserInboxRecords FindCount unreadTotal error: %v", err)
		}
	} else {
		unreadTotal, err = l.svcCtx.TNotifyRecordModel.FindCount(l.ctx,
			"channel = ? AND recipient = ? AND status = ?", "inbox", in.UserId, "unread")
		if err != nil {
			l.Errorf("ListUserInboxRecords FindCount unreadTotal error: %v", err)
		}
	}

	return &notificationrpc.ListUserInboxRecordsResponse{
		PageResult: &notificationrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		Records:     list,
		UnreadTotal: unreadTotal,
	}, nil
}

func (l *ListUserInboxRecordsLogic) getMessageMap(deliveries []*model.TNotifyRecord) map[int64]*model.TNotifyMessage {
	if len(deliveries) == 0 {
		return nil
	}

	ids := make([]interface{}, 0, len(deliveries))
	seen := make(map[int64]bool)
	for _, d := range deliveries {
		if !seen[d.MessageId] {
			ids = append(ids, d.MessageId)
			seen[d.MessageId] = true
		}
	}

	if len(ids) == 0 {
		return nil
	}

	// 构建 IN 条件
	inClause := "id IN ("
	for i := range ids {
		if i > 0 {
			inClause += ","
		}
		inClause += "?"
	}
	inClause += ")"

	msgs, err := l.svcCtx.TNotifyMessageModel.FindALL(l.ctx, inClause, ids...)
	if err != nil {
		l.Logger.Errorf("查询消息失败: %v", err)
		return nil
	}

	m := make(map[int64]*model.TNotifyMessage, len(msgs))
	for _, msg := range msgs {
		m[msg.Id] = msg
	}
	return m
}

func (l *ListUserInboxRecordsLogic) getMessageIdsByTitle(title string) []int64 {
	msgs, err := l.svcCtx.TNotifyMessageModel.FindALL(l.ctx, "title LIKE ?", "%"+title+"%")
	if err != nil {
		l.Logger.Errorf("查询消息标题失败: %v", err)
		return nil
	}
	ids := make([]int64, 0, len(msgs))
	for _, msg := range msgs {
		ids = append(ids, msg.Id)
	}
	return ids
}
