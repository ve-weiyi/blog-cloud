package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type UpdateAnnouncementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnnouncementLogic {
	return &UpdateAnnouncementLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新公告
func (l *UpdateAnnouncementLogic) UpdateAnnouncement(in *notificationrpc.UpdateAnnouncementRequest) (*notificationrpc.UpdateAnnouncementResponse, error) {
	if err := validateAudience(in.TargetType, in.TargetIds); err != nil {
		return nil, err
	}
	if in.Title == "" {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "标题不能为空")
	}

	// 存在性单独判断：值没变化时影响行数也是 0，不能拿它区分「不存在」
	if _, err := l.svcCtx.TAnnouncementModel.FindById(l.ctx, in.Id); err != nil {
		return nil, bizerr.NewBizError(bizcode.CodeResourceNotFound, "公告不存在")
	}

	// 只写可编辑字段：状态、发布时间与发布人由发布/撤回流程维护，
	// 整行覆盖会把它们一起写回旧值
	if _, err := l.svcCtx.TAnnouncementModel.UpdateFields(l.ctx, map[string]interface{}{
		"title":       in.Title,
		"content":     in.Content,
		"category":    in.Category,
		"level":       in.Level,
		"target_type": in.TargetType,
		"target_ids":  in.TargetIds,
	}, "id = ?", in.Id); err != nil {
		return nil, err
	}

	// 重新读回：响应要给调用方落库后的那一份
	entity, err := l.svcCtx.TAnnouncementModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.UpdateAnnouncementResponse{
		Success:      true,
		Announcement: convertAnnouncement(entity),
	}, nil
}
