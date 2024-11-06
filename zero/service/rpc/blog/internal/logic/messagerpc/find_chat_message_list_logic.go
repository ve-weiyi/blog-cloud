package messagerpclogic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/messagerpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindChatMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindChatMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindChatMessageListLogic {
	return &FindChatMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询聊天记录列表
func (l *FindChatMessageListLogic) FindChatMessageList(in *messagerpc.FindChatMessageListReq) (*messagerpc.FindChatMessageListResp, error) {
	page, size, sorts, conditions, params := convertChatQuery(in)

	result, err := l.svcCtx.TChatMessageModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TChatMessageModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*messagerpc.ChatMessageDetails
	for _, v := range result {
		list = append(list, convertChatMessageOut(v))
	}

	return &messagerpc.FindChatMessageListResp{
		List:  list,
		Total: count,
	}, nil
}

func convertChatQuery(in *messagerpc.FindChatMessageListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(1)
	size = int(in.Limit)

	if in.After != 0 {
		if conditions != "" {
			conditions += " and "
		}
		conditions = "created_at >= ?"
		params = append(params, time.Unix(in.After, 0))
	}

	if in.Before != 0 {
		if conditions != "" {
			conditions += " and "
		}
		conditions = "created_at <= ?"
		params = append(params, time.Unix(in.Before, 0))
	}

	if in.TopicId != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions = "topic_id = ?"
		params = append(params, in.TopicId)
	}

	if in.UserId != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions = "user_id = ?"
		params = append(params, in.UserId)
	}

	if in.Type != "" {
		if conditions != "" {
			conditions += " and "
		}
		conditions = "type = ?"
		params = append(params, in.Type)
	}

	return
}
