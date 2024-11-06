package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagLogic {
	return &AddTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建标签
func (l *AddTagLogic) AddTag(in *articlerpc.TagNewReq) (*articlerpc.TagDetails, error) {
	entity := convertTagIn(in)
	_, err := l.svcCtx.TTagModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &articlerpc.TagDetails{}, nil
}
