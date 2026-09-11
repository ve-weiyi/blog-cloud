package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
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
	status := int64(1)

	in := &socialservice.ListTalksRequest{
		PageQuery: &socialservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		Status: &status,
	}

	out, err := l.svcCtx.SocialService.ListTalks(l.ctx, in)
	if err != nil {
		return nil, err
	}

	usm, err := apiutils.BatchQuery(out.List,
		func(v *socialservice.Talk) string {
			return v.UserId
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Talk, 0)
	for _, v := range out.List {
		list = append(list, &types.Talk{
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

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
