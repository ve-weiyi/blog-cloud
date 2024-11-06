package photorpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/photorpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindBannerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindBannerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindBannerListLogic {
	return &FindBannerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询页面列表
func (l *FindBannerListLogic) FindBannerList(in *photorpc.FindBannerListReq) (*photorpc.FindBannerListResp, error) {
	page, size, sorts, conditions, params := convertBannerQuery(in)

	result, err := l.svcCtx.TBannerModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*photorpc.BannerDetails
	for _, v := range result {
		list = append(list, convertBannerOut(v))
	}

	return &photorpc.FindBannerListResp{
		List: list,
	}, nil
}

func convertBannerQuery(in *photorpc.FindBannerListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	return
}
