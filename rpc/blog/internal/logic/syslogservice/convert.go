package syslogservicelogic

import (
	"time"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

// ==================== LoginLog 转换 ====================

func convertLoginLogOut(in *model.TLoginLog) *syslogrpc.LoginLog {
	logoutAt := int64(0)
	if in.LogoutAt != nil {
		logoutAt = in.LogoutAt.UnixMilli()
	}
	return &syslogrpc.LoginLog{
		Id:         in.Id,
		UserId:     in.UserId,
		DeviceId:   in.DeviceId,
		LoginType:  in.LoginType,
		Status:     in.Status,
		FailReason: in.FailReason,
		LogoutAt:   logoutAt,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
	}
}

func convertLoginLogIn(in *syslogrpc.CreateLoginLogRequest) *model.TLoginLog {
	failReason := ""
	if in.FailReason != nil {
		failReason = *in.FailReason
	}
	return &model.TLoginLog{
		UserId:     in.UserId,
		DeviceId:   in.DeviceId,
		LoginType:  in.LoginType,
		Status:     in.Status,
		FailReason: failReason,
	}
}

// ==================== OperationLog 转换 ====================

func convertOperationLogOut(in *model.TOperationLog) *syslogrpc.OperationLog {
	return &syslogrpc.OperationLog{
		Id:             in.Id,
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		Module:         in.Module,
		Description:    in.Description,
		RequestUri:     in.RequestUri,
		RequestMethod:  in.RequestMethod,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
	}
}

func convertOperationLogIn(in *syslogrpc.CreateOperationLogRequest) *model.TOperationLog {
	return &model.TOperationLog{
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		Module:         in.Module,
		Description:    in.Description,
		RequestUri:     in.RequestUri,
		RequestMethod:  in.RequestMethod,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
	}
}

// ==================== UploadLog 转换 ====================

func convertUploadLogOut(in *model.TUploadLog) *syslogrpc.UploadLog {
	return &syslogrpc.UploadLog{
		Id:        in.Id,
		UserId:    in.UserId,
		DeviceId:  in.DeviceId,
		FileBase:  in.FileBase,
		FileName:  in.FileName,
		FileType:  in.FileType,
		FileSize:  in.FileSize,
		FileUrl:   in.FileUrl,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}

func convertUploadLogIn(in *syslogrpc.CreateUploadLogRequest) *model.TUploadLog {
	return &model.TUploadLog{
		UserId:   in.UserId,
		DeviceId: in.DeviceId,
		FileBase: in.FileBase,
		FileName: in.FileName,
		FileType: in.FileType,
		FileSize: in.FileSize,
		FileUrl:  in.FileUrl,
	}
}

// ==================== VisitLog 转换 ====================

func convertVisitLogOut(in *model.TVisitLog) *syslogrpc.VisitLog {
	return &syslogrpc.VisitLog{
		Id:        in.Id,
		UserId:    in.UserId,
		DeviceId:  in.DeviceId,
		PageName:  in.PageName,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}

func convertVisitLogIn(in *syslogrpc.CreateVisitLogRequest) *model.TVisitLog {
	return &model.TVisitLog{
		UserId:   in.UserId,
		DeviceId: in.DeviceId,
		PageName: in.PageName,
	}
}

// ==================== 时间转换辅助函数 ====================

func timeFromMilli(ms int64) *time.Time {
	if ms == 0 {
		return nil
	}
	t := time.UnixMilli(ms)
	return &t
}
