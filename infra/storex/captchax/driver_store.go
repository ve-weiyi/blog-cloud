package captchax

import (
	"context"
	"crypto/subtle"
	"time"

	"github.com/mojocn/base64Captcha"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

// driverStore 将 storex.KVStore 适配为 base64Captcha.Store。
//
// base64Captcha.Store 的接口不带 context，而生成图形验证码时由 base64Captcha
// 内部回调 Set 写入答案，调用方的 ctx 传递不进来，这里只能使用 Background。
// 校验路径不经过本适配器——Store.VerifyCaptcha 直接访问 KVStore 并传入调用方 ctx。
type driverStore struct {
	store      storex.KVStore
	keyPrefix  string
	expiration time.Duration
}

var _ base64Captcha.Store = (*driverStore)(nil)

// newDriverStore 创建适配器实例
func newDriverStore(store storex.KVStore, keyPrefix string, expiration time.Duration) *driverStore {
	return &driverStore{
		store:      store,
		keyPrefix:  keyPrefix,
		expiration: expiration,
	}
}

func (a *driverStore) Set(key string, value string) error {
	return a.store.Set(context.Background(), a.keyPrefix+key, value, a.expiration)
}

func (a *driverStore) Get(key string, clear bool) string {
	ctx := context.Background()
	val, found, err := a.store.Get(ctx, a.keyPrefix+key)
	if err != nil || !found {
		return ""
	}
	if clear {
		_ = a.store.Delete(ctx, a.keyPrefix+key)
	}
	return val
}

func (a *driverStore) Verify(key, answer string, clear bool) bool {
	v := a.Get(key, clear)
	return v != "" && subtle.ConstantTimeCompare([]byte(v), []byte(answer)) == 1
}
