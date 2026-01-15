package socialrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTalkLogic {
	return &AddTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建说说
func (l *AddTalkLogic) AddTalk(in *socialrpc.AddTalkReq) (*socialrpc.AddTalkResp, error) {
	entity := &model.TTalk{
		Id:        in.Id,
		UserId:    in.UserId,
		Content:   in.Content,
		Images:    jsonconv.AnyToJsonNE(in.ImgList),
		IsTop:     in.IsTop,
		Status:    in.Status,
		LikeCount: 0,
	}

	_, err := l.svcCtx.TTalkModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &socialrpc.AddTalkResp{
		Talk: convertTalkOut(entity),
	}, nil
}

func convertTalkOut(in *model.TTalk) (out *socialrpc.Talk) {
	var images []string
	jsonconv.JsonToAny(in.Images, &images)

	out = &socialrpc.Talk{
		Id:           in.Id,
		UserId:       in.UserId,
		Content:      in.Content,
		ImgList:      images,
		IsTop:        in.IsTop,
		Status:       in.Status,
		LikeCount:    in.LikeCount,
		CommentCount: 0,
		CreatedAt:    in.CreatedAt.UnixMilli(),
		UpdatedAt:    in.UpdatedAt.UnixMilli(),
	}

	return out
}
