package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFileLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFileLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFileLogLogic {
	return &AddFileLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文件记录
func (l *AddFileLogLogic) AddFileLog(in *syslogrpc.NewFileLogReq) (*syslogrpc.FileLogDetailsResp, error) {
	uid, _ := rpcutils.GetUserIdFromCtx(l.ctx)
	tid, _ := rpcutils.GetTerminalIdFromCtx(l.ctx)

	entity := &model.TFileLog{
		Id:         0,
		UserId:     uid,
		TerminalId: tid,
		FilePath:   in.FilePath,
		FileName:   in.FileName,
		FileType:   in.FileType,
		FileSize:   in.FileSize,
		FileMd5:    in.FileMd5,
		FileUrl:    in.FileUrl,
	}

	_, err := l.svcCtx.TFileLogModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFileLogOut(entity), nil
}
