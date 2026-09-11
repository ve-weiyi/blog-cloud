package notify_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type QueryNotifyTemplateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知模板列表
func NewQueryNotifyTemplateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryNotifyTemplateListLogic {
	return &QueryNotifyTemplateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryNotifyTemplateListLogic) QueryNotifyTemplateList(req *types.QueryNotifyTemplateListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.NotificationService.ListNotifyTemplates(l.ctx, &notificationservice.ListNotifyTemplatesRequest{
		PageQuery: &notificationservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Channel:   req.Channel,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.NotifyTemplateVO
	for _, v := range out.Templates {
		list = append(list, &types.NotifyTemplateVO{
			Id:        v.Id,
			Code:      v.Code,
			Channel:   v.Channel,
			Scene:     v.Scene,
			Title:     v.Title,
			Content:   v.Content,
			Enabled:   v.Enabled,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
