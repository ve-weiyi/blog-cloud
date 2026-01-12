package articlerpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type UpdateTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTagLogic {
	return &UpdateTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新标签
func (l *UpdateTagLogic) UpdateTag(in *articlerpc.UpdateTagReq) (*articlerpc.UpdateTagResp, error) {
	entity, err := l.svcCtx.TTagModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.TagName = in.TagName
	_, err = l.svcCtx.TTagModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &articlerpc.UpdateTagResp{
		Tag: &articlerpc.Tag{
			Id:        entity.Id,
			TagName:   entity.TagName,
			CreatedAt: entity.CreatedAt.UnixMilli(),
			UpdatedAt: entity.UpdatedAt.UnixMilli(),
		},
	}, nil
}
