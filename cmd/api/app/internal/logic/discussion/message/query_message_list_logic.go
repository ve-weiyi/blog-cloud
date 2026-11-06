package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/enrich"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/vkit/x/collectx"
)

type QueryMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取留言列表
func NewQueryMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMessageListLogic {
	return &QueryMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryMessageListLogic) QueryMessageList(req *types.QueryMessageListReq) (resp *types.ListResult, err error) {
	status := int64(enums.MessageStatusNormal)

	in := &discussionservice.ListMessagesRequest{
		ListQuery: &discussionservice.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		Status: &status,
	}

	out, err := l.svcCtx.DiscussionService.ListMessages(l.ctx, in)
	if err != nil {
		return nil, err
	}

	usm, err := enrich.GetUserInfos(l.ctx, l.svcCtx,
		collectx.Distinct(out.List, func(v *discussionservice.Message) string {
			return v.UserId
		}),
	)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Message, 0)
	for _, v := range out.List {
		list = append(list, &types.Message{
			Id:             v.Id,
			UserId:         v.UserId,
			DeviceId:       v.DeviceId,
			MessageContent: v.MessageContent,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			UserInfo:       usm[v.UserId],
		})
	}

	resp = &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}
	return
}
