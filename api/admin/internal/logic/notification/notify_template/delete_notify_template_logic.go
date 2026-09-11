package notify_template

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type DeleteNotifyTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除通知模板
func NewDeleteNotifyTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotifyTemplateLogic {
	return &DeleteNotifyTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotifyTemplateLogic) DeleteNotifyTemplate(req *types.DeleteNotifyTemplateReq) (resp *types.BatchResp, err error) {
	var successCount int64
	for _, id := range req.Ids {
		_, err = l.svcCtx.NotificationService.DeleteNotifyTemplate(l.ctx, &notificationservice.DeleteNotifyTemplateRequest{
			Id: id,
		})
		if err != nil {
			return nil, err
		}
		successCount++
	}

	return &types.BatchResp{
		SuccessCount: successCount,
	}, nil
}
