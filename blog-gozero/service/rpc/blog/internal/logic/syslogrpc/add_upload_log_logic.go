package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUploadLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUploadLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUploadLogLogic {
	return &AddUploadLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建上传记录
func (l *AddUploadLogLogic) AddUploadLog(in *syslogrpc.UploadLogNewReq) (*syslogrpc.UploadLogDetailsResp, error) {
	entity := &model.TUploadLog{
		Id:       0,
		UserId:   in.UserId,
		FilePath: in.FilePath,
		FileName: in.FileName,
		FileType: in.FileType,
		FileSize: in.FileSize,
		FileMd5:  in.FileMd5,
		FileUrl:  in.FileUrl,
		//CreatedAt:      time.Unix(in.CreatedAt, 0),
		//UpdatedAt:      time.Unix(in.UpdatedAt, 0),
	}

	_, err := l.svcCtx.TUploadLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertUploadLogOut(entity), nil
}
