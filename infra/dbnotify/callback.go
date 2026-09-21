package dbnotify

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
)

// tableChannels 定义表名到 Redis 失效频道的映射
var tableChannels = map[string][]string{
	"t_api": {
		cachekey.TraceInvalidateChannel,
		cachekey.PolicyInvalidateChannel,
	},
	"t_role": {
		cachekey.PolicyInvalidateChannel,
	},
	"t_role_api": {
		cachekey.PolicyInvalidateChannel,
	},
}

// Register 在 GORM 上注册全局回调，监听指定表的写操作并发布 Redis 失效消息
func Register(db *gorm.DB, rds *redis.Client) {
	if db == nil || rds == nil {
		return
	}

	cb := &callback{rds: rds}

	db.Callback().Create().After("gorm:create").Register("dbnotify:after_create", cb.onChange)
	db.Callback().Update().After("gorm:update").Register("dbnotify:after_update", cb.onChange)
	db.Callback().Delete().After("gorm:delete").Register("dbnotify:after_delete", cb.onChange)
}

type callback struct {
	rds *redis.Client
}

func (c *callback) onChange(db *gorm.DB) {
	if db.Statement.RowsAffected == 0 {
		return
	}

	table := db.Statement.Table
	channels, ok := tableChannels[table]
	if !ok {
		return
	}

	ctx := context.Background()
	for _, ch := range channels {
		if err := c.rds.Publish(ctx, ch, "1").Err(); err != nil {
			logx.Errorf("dbnotify publish to %s failed: %v", ch, err)
		} else {
			logx.Infof("dbnotify: table=%s -> channel=%s", table, ch)
		}
	}
}
