package talk

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建说说
func NewAddTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTalkLogic {
	return &AddTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTalkLogic) AddTalk(req *types.NewTalkReq) (resp *types.TalkBackVO, err error) {
	in := &socialrpc.NewTalkReq{
		Id:      req.Id,
		UserId:  cast.ToString(l.ctx.Value(bizheader.HeaderUid)),
		Content: req.Content,
		ImgList: req.ImgList,
		IsTop:   req.IsTop,
		Status:  req.Status,
	}

	out, err := l.svcCtx.SocialRpc.AddTalk(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertTalkTypes(out)
	return resp, nil
}

func ConvertTalkTypes(v *socialrpc.TalkDetailsResp) (out *types.TalkBackVO) {
	out = &types.TalkBackVO{
		Id:           v.Id,
		UserId:       v.UserId,
		Content:      v.Content,
		ImgList:      v.ImgList,
		IsTop:        v.IsTop,
		Status:       v.Status,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
	return
}
