package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreateTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建标签
func (l *CreateTagLogic) CreateTag(in *articlerpc.CreateTagRequest) (*articlerpc.CreateTagResponse, error) {
	entity := &model.TTag{TagName: in.TagName}
	_, err := l.svcCtx.TTagModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &articlerpc.CreateTagResponse{Id: entity.Id}, nil
}
