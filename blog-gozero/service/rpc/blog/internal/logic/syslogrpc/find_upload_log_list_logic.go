package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUploadLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUploadLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUploadLogListLogic {
	return &FindUploadLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询上传记录列表
func (l *FindUploadLogListLogic) FindUploadLogList(in *syslogrpc.FindUploadLogListReq) (*syslogrpc.FindUploadLogListResp, error) {
	page, size, sorts, conditions, params := convertUploadLogQuery(in)

	records, total, err := l.svcCtx.TUploadLogModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.UploadLogDetailsResp
	for _, v := range records {
		list = append(list, convertUploadLogOut(v))
	}

	return &syslogrpc.FindUploadLogListResp{
		List: list,
		Pagination: &syslogrpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertUploadLogQuery(in *syslogrpc.FindUploadLogListReq) (page int, size int, sorts string, conditions string, params []any) {
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

func convertUploadLogOut(in *model.TUploadLog) (out *syslogrpc.UploadLogDetailsResp) {
	out = &syslogrpc.UploadLogDetailsResp{
		Id:        in.Id,
		UserId:    in.UserId,
		FilePath:  in.FilePath,
		FileName:  in.FileName,
		FileType:  in.FileType,
		FileSize:  in.FileSize,
		FileMd5:   in.FileMd5,
		FileUrl:   in.FileUrl,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
	}

	return out
}
