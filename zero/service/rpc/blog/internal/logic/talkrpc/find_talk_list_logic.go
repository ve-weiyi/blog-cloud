package talkrpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/talkrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindTalkListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindTalkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTalkListLogic {
	return &FindTalkListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取说说列表
func (l *FindTalkListLogic) FindTalkList(in *talkrpc.FindTalkListReq) (*talkrpc.FindTalkListResp, error) {
	page, size, sorts, conditions, params := convertTalkQuery(in)

	result, err := l.svcCtx.TTalkModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*talkrpc.TalkDetails
	for _, v := range result {
		list = append(list, convertTalkOut(v))
	}

	return &talkrpc.FindTalkListResp{
		List: list,
	}, nil
}

func convertTalkQuery(in *talkrpc.FindTalkListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	if in.Status != 0 {
		conditions += " status = ?"
		params = append(params, in.Status)
	}

	return
}
