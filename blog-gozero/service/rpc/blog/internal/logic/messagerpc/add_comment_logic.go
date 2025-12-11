package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/ipx"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建评论
func (l *AddCommentLogic) AddComment(in *messagerpc.CommentNewReq) (*messagerpc.CommentDetailsResp, error) {
	uid, _ := rpcutils.GetUserIdFromCtx(l.ctx)
	ip, _ := rpcutils.GetRemoteIPFromCtx(l.ctx)
	is := ipx.GetIpSourceByBaidu(ip)

	entity := &model.TComment{
		Id:             0,
		UserId:         uid,
		TopicId:        in.TopicId,
		ParentId:       in.ParentId,
		ReplyMsgId:     in.ReplyMsgId,
		ReplyUserId:    in.ReplyUserId,
		CommentContent: in.CommentContent,
		IpAddress:      ip,
		IpSource:       is,
		Type:           in.Type,
		Status:         in.Status,
		IsReview:       l.svcCtx.Config.DefaultCommentReviewStatus,
	}

	_, err := l.svcCtx.TCommentModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertCommentOut(entity), nil
}
