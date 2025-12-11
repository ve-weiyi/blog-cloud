package online

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func Test_Online(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	sv := NewOnlineUserService(rdb, 3600)

	_ = sv.Login(ctx, "user123")
	_ = sv.Login(ctx, "user456")

	count, _ := sv.GetOnlineUserCount(ctx)
	fmt.Println("在线用户数量:", count)

	ids, _ := sv.GetOnlineUsers(ctx, 0, 0)
	fmt.Println("在线用户列表:", ids)

	_ = sv.Logout(ctx, "user123")
}
