package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesCategoryLogic {
	return &DeletesCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除文章分类
func (l *DeletesCategoryLogic) DeletesCategory(in *articlerpc.DeletesCategoryReq) (*articlerpc.DeletesCategoryResp, error) {
	rows, err := l.svcCtx.TCategoryModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &articlerpc.DeletesCategoryResp{
		SuccessCount: rows,
	}, nil
}
