package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
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

func (l *QueryMessageListLogic) QueryMessageList(req *types.QueryMessageListReq) (resp *types.PageResult, err error) {
	status := int64(1)

	in := &discussionservice.ListMessagesRequest{
		PageQuery: &discussionservice.PageQuery{
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

	usm, err := apiutils.BatchQuery(out.List,
		func(v *discussionservice.Message) string {
			return v.UserId
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
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

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
