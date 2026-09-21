package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type UpdateUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRolesLogic {
	return &UpdateUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户角色
func (l *UpdateUserRolesLogic) UpdateUserRoles(in *accessrpc.UpdateUserRolesRequest) (*accessrpc.UpdateUserRolesResponse, error) {
	// 先全删再全插必须在同一事务内：中途失败会让该用户一个角色都不剩，
	// 而这是一张授权表，失败的直接后果是权限被清空。
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TUserRoleModel.WithTx(tx).DeleteBatch(l.ctx, "user_id = ?", in.UserId); err != nil {
			return err
		}

		if len(in.RoleIds) == 0 {
			return nil
		}

		batch := make([]*model.TUserRole, 0, len(in.RoleIds))
		for _, id := range in.RoleIds {
			batch = append(batch, &model.TUserRole{
				UserId: in.UserId,
				RoleId: id,
			})
		}
		_, err := l.svcCtx.TUserRoleModel.WithTx(tx).InsertBatch(l.ctx, batch...)
		return err
	})
	if err != nil {
		return nil, err
	}

	// 通知所有实例丢弃该用户的内存角色缓存。
	//
	// 判定路径的读取顺序是「内存 -> Redis -> RPC」，且内存命中不校验年龄，
	// 因此没有这条通知时，减少角色可能长期不生效、且不产生任何报错。
	// 注意 t_user_role 不在表级失效映射里——这张表的变更不会自动触发广播，
	// 必须由变更入口自己发出。
	if err := l.publishUserRoleInvalidate(in.UserId); err != nil {
		// 权限数据已落库，不能让失效失败把变更结果一并吞掉；
		// 但必须留下痕迹——重新提交一次变更即可补发通知。
		l.Logger.Errorf("发布用户角色失效通知失败: userId=%s err=%v", in.UserId, err)
	}

	return &accessrpc.UpdateUserRolesResponse{Success: true}, nil
}

// publishUserRoleInvalidate 广播指定用户的角色失效通知，载荷即 user_id
func (l *UpdateUserRolesLogic) publishUserRoleInvalidate(userId string) error {
	if l.svcCtx.Redis == nil || userId == "" {
		return nil
	}
	return l.svcCtx.Redis.Publish(l.ctx, cachekey.UserRoleInvalidateChannel, userId).Err()
}
