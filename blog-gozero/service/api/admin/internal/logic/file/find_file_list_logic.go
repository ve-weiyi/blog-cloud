package file

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取文件列表
func NewFindFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFileListLogic {
	return &FindFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindFileListLogic) FindFileList(req *types.FileQuery) (resp *types.PageResp, err error) {
	in := &resourcerpc.FindFileUploadListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
		FilePath: req.FilePath,
		FileType: req.FileType,
	}

	out, err := l.svcCtx.ResourceRpc.FindFileUploadList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.FileBackDTO
	for _, v := range out.List {
		m := ConvertFileUploadTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertFileUploadTypes(in *resourcerpc.FileUploadDetails) (out *types.FileBackDTO) {
	out = &types.FileBackDTO{
		Id:        in.Id,
		UserId:    in.UserId,
		FilePath:  in.FilePath,
		FileName:  in.FileName,
		FileType:  in.FileType,
		FileSize:  in.FileSize,
		FileMd5:   in.FileMd5,
		FileUrl:   in.FileUrl,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
