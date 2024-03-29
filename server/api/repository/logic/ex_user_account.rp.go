package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/cache"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/constant"
)

// 注销用户
func (s *UserAccountRepository) Logoff(ctx context.Context, id int) (data any, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)
	var account entity.UserAccount

	// 开始事务
	tx := db.Begin()

	// 执行数据库操作
	db.Where("id = ?", id).First(&account)

	// 删除用户
	account.Status = constant.UserStatusDisabled
	db.Save(&account)

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		// 处理提交事务错误
		tx.Rollback() // 回滚事务
		return nil, err
	}

	return account, err
}

// 注册
func (s *UserAccountRepository) Register(ctx context.Context, account *entity.UserAccount, info *entity.UserInformation) (u *entity.UserAccount, i *entity.UserInformation, err error) {
	//var role entity.Role

	db := s.DbEngin.WithContext(ctx)
	// 开始事务
	tx := db.Begin()
	// 执行数据库操作
	/** 创建用户 **/
	db.Create(&account)

	/** 创建用户信息 start**/
	info.UserID = account.ID

	// 默认邮箱
	if info.Email == "" {
		if account.RegisterType == constant.LoginEmail {
			info.Email = account.Username
		}
	}

	// 默认昵称
	if info.Nickname == "" {
		if info.Email != "" {
			info.Nickname = strings.Split(info.Email, "@")[0]
		} else {
			info.Nickname = fmt.Sprintf("游客%v", uuid.New().String()[:8])
		}
	}

	// 默认头像
	if info.Avatar == "" {
		info.Avatar = "https://mms1.baidu.com/it/u=2815887849,1501151317&fm=253&app=138&f=JPEG"
	}

	// 默认简介
	if info.Intro == "" {
		info.Intro = "这个人很神秘，什么都没有写！"
	}

	db.Create(&info)
	/** 创建用户信息 end **/

	/** 创建用户角色 end **/
	var roles []*entity.Role
	db.Where("is_default = ?", 1).Find(&roles)

	var userRoles []*entity.UserRole
	for _, item := range roles {
		userRoles = append(userRoles, &entity.UserRole{
			UserID: account.ID,
			RoleID: item.ID,
		})
	}

	db.Create(&userRoles)
	/** 创建用户角色 end **/

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		// 处理提交事务错误
		tx.Rollback() // 回滚事务
		return nil, nil, err
	}

	return account, info, nil
}

func (s *UserAccountRepository) FindUserMenus(ctx context.Context, userId int) (list []*entity.Menu, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)

	//查询用户信息
	var account entity.UserAccount
	err = db.Where("user_id = ?", userId).First(&account).Error
	if err != nil {
		return nil, err
	}

	//查询用户的角色
	var roleApis []*entity.RoleMenu
	err = db.Where("user_id = ?", userId).Find(&roleApis).Error
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, item := range roleApis {
		ids = append(ids, item.MenuID)
	}

	var apis []*entity.Menu
	err = db.Where("id in (?)", ids).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	return apis, nil
}

// 加载用户model
func (s *UserAccountRepository) LoadUserByUsername(ctx context.Context, username string) (data *entity.UserAccount, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)

	//查询用户信息
	err = db.Where("username = ?", username).First(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 根据条件获取UserOauth记录
func (s *UserAccountRepository) FindUserOauthByOpenid(ctx context.Context, openId string, platform string) (out *entity.UserOauth, err error) {
	db := s.DbEngin.WithContext(ctx)
	err = db.Where("open_id = ? and platform = ?", openId, platform).First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 查询用户信息
func (s *UserAccountRepository) FindUserInfo(ctx context.Context, userId int) (out *entity.UserInformation, err error) {
	// 创建db
	db := s.DbEngin.WithContext(ctx)

	//查询用户信息
	err = db.Where("user_id = ?", userId).First(&out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

// 根据id获取UserLoginHistory记录
func (s *UserAccountRepository) FindLastLoginHistory(ctx context.Context, uid int) (out *entity.UserLoginHistory, err error) {
	db := s.DbEngin.WithContext(ctx)

	out = &entity.UserLoginHistory{}
	err = db.Where("user_id = ?", uid).First(&out).Error
	if err != nil {
		return nil, err
	}
	return out, err
}

// 用户登录
func (s *UserAccountRepository) Login(ctx context.Context, account *entity.UserAccount) (data string, err error) {
	// redis 保存用户登录信息
	onlineKey := cache.WrapCacheKey(cache.UserOnline, account.ID)
	return s.Cache.Set(ctx, onlineKey, account.Username, cache.ExpireTimeWeek).Result()
}

// 用户登出
func (s *UserAccountRepository) Logout(ctx context.Context, uid int) (data int64, err error) {
	// redis 删除用户登录信息
	onlineKey := cache.WrapCacheKey(cache.UserOnline, uid)
	return s.Cache.Del(ctx, onlineKey).Result()
}

// 在线用户
func (s *UserAccountRepository) Online(ctx context.Context, page int, pageSize int) (data []string, err error) {
	redisClient := s.Cache
	// redis 获取在线用户
	var cursor uint64
	var onlineUsers []string

	// 使用 SCAN 迭代遍历键空间
	for {
		// 使用 SCAN 命令获取部分键
		keys, newCursor, err := redisClient.Scan(ctx, cursor, cache.WrapCacheKey(cache.UserOnline, "*"), int64(pageSize)).Result()
		if err != nil {
			return nil, err
		}

		// 更新游标
		cursor = newCursor

		// 直接添加到在线用户列表
		onlineUsers = append(onlineUsers, keys...)

		// 如果游标为 0，则说明迭代结束
		if cursor == 0 {
			break
		}
	}

	// 实现分页逻辑
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex > len(onlineUsers) {
		endIndex = len(onlineUsers)
	}

	for _, value := range onlineUsers[startIndex:endIndex] {
		data = append(data, strings.TrimPrefix(value, cache.UserOnline+":"))
	}

	return data, nil
}
