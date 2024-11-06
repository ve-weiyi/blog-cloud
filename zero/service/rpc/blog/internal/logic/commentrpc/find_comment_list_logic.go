package commentrpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentListLogic {
	return &FindCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取评论列表
func (l *FindCommentListLogic) FindCommentList(in *commentrpc.FindCommentListReq) (*commentrpc.FindCommentListResp, error) {
	page, size, sorts, conditions, params := convertCommentQuery(in)

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

	return &commentrpc.FindCommentListResp{
		List:  list,
		Total: count,
	}, nil
}

func convertCommentQuery(in *commentrpc.FindCommentListReq) (page int, size int, sorts string, conditions string, params []any) {
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
