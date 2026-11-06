package aiservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/airpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

const (
	defaultSessionPageSize = 20
	maxSessionPageSize     = 100
)

type ListSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionsLogic {
	return &ListSessionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询会话列表
func (l *ListSessionsLogic) ListSessions(in *airpc.ListSessionsRequest) (*airpc.ListSessionsResponse, error) {
	uid, err := requireUserID(l.ctx)
	if err != nil {
		return nil, err
	}

	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	size := int(in.PageSize)
	if size <= 0 {
		size = defaultSessionPageSize
	}
	if size > maxSessionPageSize {
		size = maxSessionPageSize
	}

	records, total, err := l.svcCtx.TAiSessionModel.FindListAndTotal(l.ctx, page, size,
		"updated_at desc, id desc", "user_id = ?", uid)
	if err != nil {
		return nil, err
	}

	list := make([]*airpc.Session, 0, len(records))
	for _, record := range records {
		list = append(list, convertSession(record))
	}

	return &airpc.ListSessionsResponse{
		ListResult: &airpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
