package cachekey

import "fmt"

// Permission模块缓存Key定义

const (
	PolicyInvalidateChannel   = "blog:permissionx:invalidate:policy"
	UserRoleInvalidateChannel = "blog:permissionx:invalidate:user"
	userRoleCachePrefix       = "blog:permissionx:user_roles:"
)

func UserRoleCacheKey(userId string) string {
	return fmt.Sprintf("%s%s", userRoleCachePrefix, userId)
}
