package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
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
	out, err := l.svcCtx.DiscussionService.ListMessages(l.ctx, &discussionservice.ListMessagesRequest{
		PageQuery: &discussionservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		UserId:    req.UserId,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}

	var uids, tids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
		tids = append(tids, v.DeviceId)
	}

	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	vsm, err := apiutils.GetGuests(l.ctx, l.svcCtx, tids)
	if err != nil {
		return nil, err
	}

	var list []*types.MessageVO
	for _, v := range out.List {
		list = append(list, &types.MessageVO{
			Id:             v.Id,
			UserId:         v.UserId,
			DeviceId:       v.DeviceId,
			MessageContent: v.MessageContent,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			UserInfo:       usm[v.UserId],
			GuestInfo:      vsm[v.DeviceId],
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
