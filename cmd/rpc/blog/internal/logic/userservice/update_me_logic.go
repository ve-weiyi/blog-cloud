package userservicelogic

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type userInfo struct {
	Gender  int64  `json:"gender"`
	Intro   string `json:"intro"`
	Website string `json:"website"`
}

type UpdateMeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMeLogic {
	return &UpdateMeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新当前用户基础信息
func (l *UpdateMeLogic) UpdateMe(in *userrpc.UpdateMeRequest) (*userrpc.UpdateMeResponse, error) {
	// 从上下文获取用户ID
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	updateFields := make(map[string]interface{})
	if in.Nickname != nil {
		updateFields["nickname"] = *in.Nickname
	}
	if in.Avatar != nil {
		updateFields["avatar"] = *in.Avatar
	}

	if in.Gender != nil || in.Intro != nil || in.Website != nil {
		user, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, userID)
		if err != nil {
			return nil, err
		}

		info := userInfo{}
		if user.Info != nil {
			if err := json.Unmarshal([]byte(*user.Info), &info); err != nil {
				return nil, err
			}
		}
		if in.Gender != nil {
			info.Gender = *in.Gender
		}
		if in.Intro != nil {
			info.Intro = *in.Intro
		}
		if in.Website != nil {
			info.Website = *in.Website
		}

		payload, err := json.Marshal(info)
		if err != nil {
			return nil, err
		}
		updateFields["info"] = string(payload)
	}

	if len(updateFields) == 0 {
		return &userrpc.UpdateMeResponse{}, nil
	}

	if _, err = l.svcCtx.TUserModel.UpdateFields(l.ctx, updateFields, "user_id = ?", userID); err != nil {
		return nil, err
	}

	return &userrpc.UpdateMeResponse{}, nil
}
