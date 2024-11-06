package tag

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewAddTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagLogic {
	return &AddTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTagLogic) AddTag(req *types.TagNewReq) (resp *types.TagBackDTO, err error) {
	in := ConvertTagPb(req)
	out, err := l.svcCtx.ArticleRpc.AddTag(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertTagTypes(out)
	return resp, nil
}

func ConvertTagPb(in *types.TagNewReq) (out *articlerpc.TagNewReq) {
	out = &articlerpc.TagNewReq{
		Id:      in.Id,
		TagName: in.TagName,
	}

	return
}

func ConvertTagTypes(in *articlerpc.TagDetails) (out *types.TagBackDTO) {
	out = &types.TagBackDTO{
		Id:           in.Id,
		TagName:      in.TagName,
		ArticleCount: in.ArticleCount,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}

	return
}
