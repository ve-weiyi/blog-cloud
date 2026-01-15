package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindFileLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindFileLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFileLogListLogic {
	return &FindFileLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文件记录列表
func (l *FindFileLogListLogic) FindFileLogList(in *syslogrpc.FindFileLogListReq) (*syslogrpc.FindFileLogListResp, error) {
	page, size, sorts, conditions, params := convertFileLogQuery(in)

	records, total, err := l.svcCtx.TFileLogModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.FileLog
	for _, v := range records {
		list = append(list, convertFileLogOut(v))
	}

	return &syslogrpc.FindFileLogListResp{
		List: list,
		Pagination: &syslogrpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertFileLogQuery(in *syslogrpc.FindFileLogListReq) (page int, size int, sorts string, conditions string, params []any) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.FilePath != "" {
		opts = append(opts, query.WithCondition("file_path like ?", in.FilePath+"%"))
	}

	if in.FileName != "" {
		opts = append(opts, query.WithCondition("file_name like ?", in.FileName+"%"))
	}

	if in.FileType != "" {
		opts = append(opts, query.WithCondition("file_type = ?", in.FileType))
	}

	return query.NewQueryBuilder(opts...).Build()
}

func convertFileLogOut(in *model.TFileLog) (out *syslogrpc.FileLog) {
	out = &syslogrpc.FileLog{
		Id:         in.Id,
		UserId:     in.UserId,
		TerminalId: in.TerminalId,
		FilePath:   in.FilePath,
		FileName:   in.FileName,
		FileType:   in.FileType,
		FileSize:   in.FileSize,
		FileMd5:    in.FileMd5,
		FileUrl:    in.FileUrl,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
	}

	return out
}
