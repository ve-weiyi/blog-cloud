package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/socialservice"
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

func (l *QueryTalkListLogic) QueryTalkList(req *types.QueryTalkListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.SocialService.ListTalks(l.ctx, &socialservice.ListTalksRequest{
		PageQuery: &socialservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
	}

	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	var list []*types.TalkVO
	for _, v := range out.List {
		list = append(list, &types.TalkVO{
			Id:           v.Id,
			UserId:       v.UserId,
			Content:      v.Content,
			ImgList:      v.Images,
			IsTop:        v.IsTop,
			Status:       v.Status,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
			UserInfo:     usm[v.UserId],
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
