package notify_record

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type QueryUserInboxRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户 inbox 投递列表
func NewQueryUserInboxRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserInboxRecordListLogic {
	return &QueryUserInboxRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserInboxRecordListLogic) QueryUserInboxRecordList(req *types.QueryUserInboxRecordListReq) (resp *types.QueryUserInboxRecordListResp, err error) {
	out, err := l.svcCtx.NotificationService.ListUserInboxRecords(l.ctx, &notificationservice.ListUserInboxRecordsRequest{
		ListQuery:  &notificationservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		UserId:     req.UserId,
		OnlyUnread: req.OnlyUnread,
		Title:      req.Title,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.NotifyRecordVO
	for _, v := range out.List {
		list = append(list, &types.NotifyRecordVO{
			Id:           v.Id,
			MessageId:    v.MessageId,
			Channel:      v.Channel,
			Recipient:    v.Recipient,
			TemplateCode: v.TemplateCode,
			Content:      v.Content,
			Status:       v.Status,
			BizId:        v.BizId,
			ErrorMsg:     v.ErrorMsg,
			ReadAt:       v.ReadAt,
			SentAt:       v.SentAt,
			CreatedAt:    v.CreatedAt,
			Title:        v.Title,
			Category:     v.Category,
			Level:        v.Level,
		})
	}

	return &types.QueryUserInboxRecordListResp{
		Page:        out.ListResult.Page,
		PageSize:    out.ListResult.PageSize,
		Total:       out.ListResult.Total,
		UnreadTotal: out.UnreadTotal,
		List:        list,
	}, nil
}
