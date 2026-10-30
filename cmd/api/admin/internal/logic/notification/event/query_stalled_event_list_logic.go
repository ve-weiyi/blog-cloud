package event

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type QueryStalledEventListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询停滞的通知事件
func NewQueryStalledEventListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryStalledEventListLogic {
	return &QueryStalledEventListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryStalledEventListLogic) QueryStalledEventList(req *types.QueryStalledEventListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.NotificationService.ListStalledEvents(l.ctx, &notificationservice.ListStalledEventsRequest{
		BeforeSeconds: derefInt64(req.BeforeSeconds),
		Limit:         derefInt64(req.Limit),
	})
	if err != nil {
		return nil, err
	}

	list := make([]*types.NotificationEventVO, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, convertEvent(item))
	}

	// 该接口不参与分页，一次最多返回 limit 条，总数即列表长度
	return &types.ListResult{
		Page:     1,
		PageSize: int64(len(list)),
		Total:    int64(len(list)),
		List:     list,
	}, nil
}

func derefInt64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}
