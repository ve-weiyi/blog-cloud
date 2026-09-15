package upload

import (
	"context"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type CreateUploadTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取上传凭证（前端直传）
func NewCreateUploadTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUploadTokenLogic {
	return &CreateUploadTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUploadTokenLogic) CreateUploadToken(req *types.CreateUploadTokenReq) (resp *types.CreateUploadTokenResp, err error) {
	// 默认过期时间为1小时
	expireSeconds := req.ExpireSeconds
	if expireSeconds <= 0 {
		expireSeconds = 3600
	}

	// 获取上传凭证
	token, err := l.svcCtx.StorageProvider.GetUploadToken(l.ctx, filepath.Join(req.FileBase, req.FileName), expireSeconds)
	if err != nil {
		l.Errorf("Failed to get upload token: %v", err)
		return nil, err
	}

	// 转换为响应格式
	resp = &types.CreateUploadTokenResp{
		UploadUrl: token.UploadURL,
		Token:     token.Token,
		Policy:    token.Policy,
		Signature: token.Signature,
		FileKey:   token.FileKey,
		AccessUrl: token.AccessURL,
		ExpireAt:  token.ExpireAt.UnixMilli(),
		ExtraData: token.ExtraData,
	}

	// 记录上传日志
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)
	did, _ := metax.GetApiDeviceIdFromCtx(l.ctx)
	_, _ = l.svcCtx.SyslogService.CreateUploadLog(l.ctx, &syslogservice.CreateUploadLogRequest{
		UserId:   uid,
		DeviceId: did,
		FileBase: req.FileBase,
		FileName: req.FileName,
		FileType: filepath.Ext(req.FileName),
		FileUrl:  token.AccessURL,
	})

	l.Infof("[GetUploadToken] Token generated for file: %s, expires at: %v", req.FileName, token.ExpireAt)
	return resp, nil
}
