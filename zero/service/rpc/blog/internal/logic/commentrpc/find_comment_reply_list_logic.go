package commentrpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCommentReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindCommentReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentReplyListLogic {
	return &FindCommentReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询评论回复列表
func (l *FindCommentReplyListLogic) FindCommentReplyList(in *commentrpc.FindCommentReplyListReq) (*commentrpc.FindCommentReplyListResp, error) {
	page, size, sorts, conditions, params := convertCommentReplyQuery(in)

	result, err := l.svcCtx.TCommentModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TCommentModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*commentrpc.CommentDetails
	for _, v := range result {
		m := convertCommentOut(v)
		list = append(list, m)
	}

	return &commentrpc.FindCommentReplyListResp{
		List:  list,
		Total: count,
	}, nil
}

func convertCommentReplyQuery(in *commentrpc.FindCommentReplyListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	if in.Type != 0 {
		conditions += " type = ?"
		params = append(params, in.Type)
	}

	if in.TopicId != 0 {
		if conditions != "" {
			conditions += " and "
		}
		conditions += " topic_id = ?"
		params = append(params, in.TopicId)
	}

	if conditions != "" {
		conditions += " and "
	}
	conditions += " parent_id = ?"
	params = append(params, in.ParentId)

	return
}
