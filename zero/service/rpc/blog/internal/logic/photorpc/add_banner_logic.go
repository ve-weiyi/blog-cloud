package photorpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/photorpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBannerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBannerLogic {
	return &AddBannerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建页面
func (l *AddBannerLogic) AddBanner(in *photorpc.BannerNewReq) (*photorpc.BannerDetails, error) {
	entity := convertBannerIn(in)

	_, err := l.svcCtx.TBannerModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertBannerOut(entity), nil
}

func convertBannerIn(in *photorpc.BannerNewReq) (out *model.TBanner) {
	out = &model.TBanner{
		Id:          in.Id,
		BannerName:  in.BannerName,
		BannerLabel: in.BannerLabel,
		BannerCover: in.BannerCover,
	}

	return out

}

func convertBannerOut(in *model.TBanner) (out *photorpc.BannerDetails) {
	out = &photorpc.BannerDetails{
		Id:          in.Id,
		BannerName:  in.BannerName,
		BannerLabel: in.BannerLabel,
		BannerCover: in.BannerCover,
		CreatedAt:   in.CreatedAt.Unix(),
		UpdatedAt:   in.UpdatedAt.Unix(),
	}

	return out
}
