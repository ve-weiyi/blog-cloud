package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/enrich"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
)

type QueryTalkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取说说列表
func NewQueryTalkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTalkListLogic {
	return &QueryTalkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTalkListLogic) QueryTalkList(req *types.QueryTalkListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.DiscussionService.ListTalks(l.ctx, &discussionservice.ListTalksRequest{
		ListQuery: &discussionservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
	}

	usm, err := enrich.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	var list []*types.TalkVO
	for _, v := range out.List {
		list = append(list, convertTalkVO(v, usm))
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
