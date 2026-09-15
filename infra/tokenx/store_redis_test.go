package tokenx_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

func setupRedisStore(t *testing.T) (tokenx2.Store, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return tokenx2.NewRedisStore(client), mr
}

func TestRedisStore_SetAndGet(t *testing.T) {
	store, _ := setupRedisStore(t)
	ctx := context.Background()

	err := store.Set(ctx, "key1", "value1", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	val, err := store.Get(ctx, "key1")
	if err != nil {
		t.Fatal(err)
	}
	if val != "value1" {
		t.Fatalf("expected value1, got %s", val)
	}
}

func TestRedisStore_Get_MissingKey_ReturnsEmptyString(t *testing.T) {
	store, _ := setupRedisStore(t)
	ctx := context.Background()

	val, err := store.Get(ctx, "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Fatalf("expected empty string for missing key, got %s", val)
	}
}

func TestRedisStore_Get_ExpiredKey_ReturnsEmptyString(t *testing.T) {
	store, mr := setupRedisStore(t)
	ctx := context.Background()

	store.Set(ctx, "key1", "value1", time.Second)
	mr.FastForward(2 * time.Second) // expire all keys

	val, err := store.Get(ctx, "key1")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Fatalf("expected empty string for expired key, got %s", val)
	}
}

func TestRedisStore_Delete(t *testing.T) {
	store, _ := setupRedisStore(t)
	ctx := context.Background()

	store.Set(ctx, "k1", "v1", time.Hour)
	store.Set(ctx, "k2", "v2", time.Hour)

	err := store.Delete(ctx, "k1", "k2")
	if err != nil {
		t.Fatal(err)
	}

	v1, _ := store.Get(ctx, "k1")
	v2, _ := store.Get(ctx, "k2")
	if v1 != "" || v2 != "" {
		t.Fatal("expected both keys deleted")
	}
}

func TestRedisStore_Keys(t *testing.T) {
	store, _ := setupRedisStore(t)
	ctx := context.Background()

	store.Set(ctx, "app:access:u1:d1", "t1", time.Hour)
	store.Set(ctx, "app:access:u1:d2", "t2", time.Hour)
	store.Set(ctx, "app:refresh:u2:d1", "t3", time.Hour)

	keys, err := store.Keys(ctx, "app:access:*")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
}

func TestRedisStore_Keys_NoMatch(t *testing.T) {
	store, _ := setupRedisStore(t)
	ctx := context.Background()

	store.Set(ctx, "app:access:u1:d1", "t1", time.Hour)

	keys, err := store.Keys(ctx, "nonexistent:*")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(keys))
	}
}
