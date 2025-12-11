package messagerpclogic

import (
	"context"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindChatListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindChatListLogic {
	return &FindChatListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询聊天记录列表
func (l *FindChatListLogic) FindChatList(in *messagerpc.FindChatListReq) (*messagerpc.FindChatListResp, error) {
	page, size, sorts, conditions, params := convertChatQuery(in)

	records, total, err := l.svcCtx.TChatModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*messagerpc.ChatDetailsResp
	for _, v := range records {
		list = append(list, convertChatOut(v))
	}

	return &messagerpc.FindChatListResp{
		List: list,
		Pagination: &messagerpc.PageResp{
			Total: total,
		},
	}, nil
}

func convertChatQuery(in *messagerpc.FindChatListReq) (page int, size int, sorts string, conditions string, params []any) {
	opts := []query.Option{
		query.WithPage(1),
		query.WithSize(int(in.Limit)),
		query.WithSorts("created_at desc"),
	}

	if in.After != 0 {
		opts = append(opts, query.WithCondition("created_at >= ?", time.Unix(in.After, 0)))
	}

	if in.Before != 0 {
		opts = append(opts, query.WithCondition("created_at <= ?", time.Unix(in.Before, 0)))
	}

	if in.UserId != "" {
		opts = append(opts, query.WithCondition("user_id = ?", in.UserId))
	}

	if in.Type != "" {
		opts = append(opts, query.WithCondition("type = ?", in.Type))
	}

	if in.Status != 0 {
		opts = append(opts, query.WithCondition("status = ?", in.Status))
	}

	return query.NewQueryBuilder(opts...).Build()
}
