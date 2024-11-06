package resourcerpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindFileUploadListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindFileUploadListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFileUploadListLogic {
	return &FindFileUploadListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文件上传列表
func (l *FindFileUploadListLogic) FindFileUploadList(in *resourcerpc.FindFileUploadListReq) (*resourcerpc.FindFileUploadListResp, error) {
	page, size, sorts, conditions, params := convertFileUploadQuery(in)

	result, err := l.svcCtx.TFileUploadModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TFileUploadModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*resourcerpc.FileUploadDetails
	for _, v := range result {
		list = append(list, convertFileUploadOut(v))
	}

	return &resourcerpc.FindFileUploadListResp{
		List:  list,
		Total: count,
	}, nil
}

func convertFileUploadQuery(in *resourcerpc.FindFileUploadListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	if in.FilePath != "" {
		conditions += " file_path = ?"
		params = append(params, in.FilePath)
	}

	if in.FileType != "" {
		conditions += " file_type = ?"
		params = append(params, in.FileType)
	}

	return
}
