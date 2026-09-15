package notify_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type BatchDeleteNotifyTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除通知模板
func NewBatchDeleteNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteNotifyTemplateLogic {
	return &BatchDeleteNotifyTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteNotifyTemplateLogic) BatchDeleteNotifyTemplate(req *types.DeleteNotifyTemplateReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.BatchDeleteNotifyTemplates(l.ctx, &notificationservice.BatchDeleteNotifyTemplatesRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
