package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindRemarkListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindRemarkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRemarkListLogic {
	return &FindRemarkListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询留言列表
func (l *FindRemarkListLogic) FindRemarkList(in *messagerpc.FindRemarkListReq) (*messagerpc.FindRemarkListResp, error) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.UserId != "" {
		opts = append(opts, query.WithCondition("user_id = ?", in.UserId))
	}

	if in.Status >= 0 {
		opts = append(opts, query.WithCondition("status = ?", in.Status))
	}

	page, size, sorts, conditions, params := query.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TRemarkModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*messagerpc.Remark
	for _, v := range records {
		list = append(list, convertRemarkOut(v))
	}

	return &messagerpc.FindRemarkListResp{
		List: list,
		Pagination: &messagerpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertRemarkOut(in *model.TRemark) (out *messagerpc.Remark) {
	out = &messagerpc.Remark{
		Id:             in.Id,
		UserId:         in.UserId,
		TerminalId:     in.TerminalId,
		MessageContent: in.MessageContent,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
	}

	return out
}
