package talkrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/talkrpc"
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
func (l *AddTalkLogic) AddTalk(in *talkrpc.TalkNewReq) (*talkrpc.TalkDetails, error) {
	entity := convertTalkIn(in)

	_, err := l.svcCtx.TTalkModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertTalkOut(entity), nil
}

func convertTalkIn(in *talkrpc.TalkNewReq) (out *model.TTalk) {
	out = &model.TTalk{
		Id:        in.Id,
		UserId:    in.UserId,
		Content:   in.Content,
		Images:    jsonconv.AnyToJsonNE(in.ImgList),
		IsTop:     in.IsTop,
		Status:    in.Status,
		LikeCount: 0,
	}

	return out
}

func convertTalkOut(in *model.TTalk) (out *talkrpc.TalkDetails) {
	var images []string
	jsonconv.JsonToAny(in.Images, &images)

	out = &talkrpc.TalkDetails{
		Id:           in.Id,
		UserId:       in.UserId,
		Content:      in.Content,
		ImgList:      images,
		IsTop:        in.IsTop,
		Status:       in.Status,
		LikeCount:    in.LikeCount,
		CommentCount: 0,
		CreatedAt:    in.CreatedAt.Unix(),
		UpdatedAt:    in.UpdatedAt.Unix(),
	}

	return out
}
