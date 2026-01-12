package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
)

type GetTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询说说
func NewGetTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTalkLogic {
	return &GetTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTalkLogic) GetTalk(req *types.IdReq) (resp *types.Talk, err error) {
	in := &socialrpc.GetTalkReq{
		Id: req.Id,
	}

	out, err := l.svcCtx.SocialRpc.GetTalk(l.ctx, in)
	if err != nil {
		return nil, err
	}

	uids := []string{out.Talk.UserId}
	// 查询用户信息
	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	// 查询评论量
	counts, err := l.svcCtx.MessageRpc.FindCommentReplyCounts(l.ctx, &messagerpc.FindCommentReplyCountsReq{
		TopicIds: []int64{out.Talk.Id},
	})
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.SyslogRpc.AddVisitLog(l.ctx, &syslogrpc.AddVisitLogReq{
		PageName: "说说",
	})
	if err != nil {
		return nil, err
	}

	resp = convertTalkTypes(out.Talk, usm, counts.TopicCommentCounts)
	return resp, nil
}
