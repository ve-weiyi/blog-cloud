package chatservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListChatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListChatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListChatsLogic {
	return &ListChatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListChatsLogic) ListChats(in *chatrpc.ListChatsRequest) (*chatrpc.ListChatsResponse, error) {
	var opts []queryx.Option
	if in.After > 0 {
		opts = append(opts, queryx.WithCondition("created_at > ?", in.After))
	}
	if in.Before > 0 {
		opts = append(opts, queryx.WithCondition("created_at < ?", in.Before))
	}
	if in.UserId != nil {
		opts = append(opts, queryx.WithCondition("user_id = ?", *in.UserId))
	}
	if in.Type != nil {
		opts = append(opts, queryx.WithCondition("type = ?", *in.Type))
	}
	if in.Content != nil {
		opts = append(opts, queryx.WithCondition("content like ?", "%"+*in.Content+"%"))
	}
	if in.Status != nil {
		opts = append(opts, queryx.WithCondition("status = ?", *in.Status))
	}

	page := 1
	size := int(in.Limit)
	if size <= 0 {
		size = 20
	}
	opts = append(opts, queryx.WithPage(page))
	opts = append(opts, queryx.WithSize(size))
	opts = append(opts, queryx.WithSorts("created_at desc"))

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TChatModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*chatrpc.Chat
	for _, entity := range records {
		list = append(list, convertChatOut(entity))
	}

	return &chatrpc.ListChatsResponse{
		ListResult: &chatrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
