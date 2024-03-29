package initialize

import (
	"time"

	"github.com/orca-zhang/ecache"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
)

func OtherInit() {
	global.BlackCache = ecache.NewLRUCache(16, 200, 10*time.Second).LRU2(1024)
}
