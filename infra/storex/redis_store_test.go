package storex_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

// newTestRedis 起一个 miniredis，返回连上它的客户端与实例本身。
// 返回实例是因为部分用例要 FastForward 或直接查看底层 key。
func newTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()

	mr := miniredis.RunT(t)

	return redis.NewClient(&redis.Options{Addr: mr.Addr()}), mr
}

func TestRedisStore_SetGet(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisStore(client)

	if err := s.Set(ctx, "k", "v", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}

	val, found, err := s.Get(ctx, "k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !found {
		t.Error("found = false, want true")
	}
	if val != "v" {
		t.Errorf("val = %q, want %q", val, "v")
	}
}

// found 区分「没写过」与「写了个空串」——
// 两者都表现为空字符串，光看 value 分不出来。
func TestRedisStore_Get_DistinguishesMissingFromEmpty(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisStore(client)

	if err := s.Set(ctx, "empty", "", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}

	tests := []struct {
		name      string
		key       string
		wantVal   string
		wantFound bool
	}{
		{name: "没写过的 key", key: "never-set"},
		{name: "写过空串的 key", key: "empty", wantFound: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, found, err := s.Get(ctx, tc.key)
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			if found != tc.wantFound {
				t.Errorf("found = %v, want %v", found, tc.wantFound)
			}
			if val != tc.wantVal {
				t.Errorf("val = %q, want %q", val, tc.wantVal)
			}
		})
	}
}

func TestRedisStore_Get_Expired(t *testing.T) {
	ctx := context.Background()
	client, mr := newTestRedis(t)
	s := storex.NewRedisStore(client)

	if err := s.Set(ctx, "k", "v", time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}
	mr.FastForward(2 * time.Second)

	val, found, err := s.Get(ctx, "k")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if found {
		t.Errorf("过期 key 仍 found = true, val = %q", val)
	}
}

// ttl <= 0 表示不过期。
func TestRedisStore_Set_NonPositiveTTL_NeverExpires(t *testing.T) {
	tests := []struct {
		name string
		ttl  time.Duration
	}{
		{name: "ttl 为 0", ttl: 0},
		{name: "ttl 为负", ttl: -time.Hour},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client, mr := newTestRedis(t)
			s := storex.NewRedisStore(client)

			if err := s.Set(ctx, "k", "v", tc.ttl); err != nil {
				t.Fatalf("Set: %v", err)
			}

			mr.FastForward(365 * 24 * time.Hour)

			if _, found, _ := s.Get(ctx, "k"); !found {
				t.Error("key 过期了，want 永不过期")
			}
			// miniredis 在未设置过期时返回 0（真实 Redis 返回 -1）
			if ttl := mr.TTL("k"); ttl != 0 {
				t.Errorf("底层 TTL = %v, want 0（未设置过期）", ttl)
			}
		})
	}
}

func TestRedisStore_Delete(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisStore(client)

	if err := s.Set(ctx, "k", "v", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Delete(ctx, "k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if _, found, _ := s.Get(ctx, "k"); found {
		t.Error("删除后仍能读到")
	}
}

// Delete 支持一次删多个 key —— tokenx 的 Revoke 要一次删掉
// access / refresh / session 三个，单 key 版本会把它拆成 3 次 RTT。
func TestRedisStore_Delete_MultipleKeys(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisStore(client)

	keys := []string{"k1", "k2", "k3"}
	for _, k := range keys {
		if err := s.Set(ctx, k, "v", time.Hour); err != nil {
			t.Fatalf("Set(%s): %v", k, err)
		}
	}

	if err := s.Delete(ctx, keys[0], keys[1:]...); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	for _, k := range keys {
		if _, found, _ := s.Get(ctx, k); found {
			t.Errorf("%s 删除后仍能读到", k)
		}
	}
}

func TestRedisStore_Delete_MissingKeyIsNotAnError(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisStore(client)

	if err := s.Delete(ctx, "never-set"); err != nil {
		t.Errorf("删不存在的 key 报错：%v", err)
	}
}

func TestRedisStore_Keys(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisStore(client)

	for _, k := range []string{"app:access:u1:d1", "app:access:u1:d2", "app:refresh:u2:d1"} {
		if err := s.Set(ctx, k, "v", time.Hour); err != nil {
			t.Fatalf("Set(%s): %v", k, err)
		}
	}

	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "按 pattern 匹配多个 key",
			pattern: "app:access:*",
			want:    []string{"app:access:u1:d1", "app:access:u1:d2"},
		},
		{name: "无匹配时返回空集合", pattern: "nonexistent:*"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Keys(ctx, tc.pattern)
			if err != nil {
				t.Fatalf("Keys: %v", err)
			}
			// SCAN 不保证顺序
			slices.Sort(got)
			if !slices.Equal(got, tc.want) {
				t.Errorf("Keys(%q) = %v, want %v", tc.pattern, got, tc.want)
			}
		})
	}
}
