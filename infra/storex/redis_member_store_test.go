package storex_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

// seedMember 直接把过期时间戳写进 ZSET，绕开 Add 的 TTL 计算。
// 这样「过期 member 的读取行为」与被测的 Add 收到什么 TTL 无关，
// 将来再动 Add 的 TTL 语义时这组用例不受影响。
func seedMember(t *testing.T, mr *miniredis.Miniredis, key, member string, expireAt time.Time) {
	t.Helper()

	if _, err := mr.ZAdd(key, float64(expireAt.Unix()), member); err != nil {
		t.Fatalf("ZAdd(%s, %s): %v", key, member, err)
	}
}

func TestRedisMemberStore_AddMembers(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	for _, m := range []string{"d1", "d2"} {
		if err := s.Add(ctx, "u1", m, time.Hour); err != nil {
			t.Fatalf("Add(%s): %v", m, err)
		}
	}

	got, err := s.Members(ctx, "u1")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if want := []string{"d1", "d2"}; !slices.Equal(got, want) {
		t.Errorf("Members = %v, want %v", got, want)
	}
}

// Members 对不存在的 key 返回空集合，不报错。
func TestRedisMemberStore_Members_MissingKey(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	got, err := s.Members(ctx, "never-set")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Members = %v, want 空", got)
	}
}

func TestRedisMemberStore_Has(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	if err := s.Add(ctx, "u1", "d1", time.Hour); err != nil {
		t.Fatalf("Add: %v", err)
	}

	tests := []struct {
		name   string
		member string
		want   bool
	}{
		{name: "已加入的 member", member: "d1", want: true},
		{name: "没加入的 member", member: "d9"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := s.Has(ctx, "u1", tc.member)
			if err != nil {
				t.Fatalf("Has: %v", err)
			}
			if got != tc.want {
				t.Errorf("Has(%s) = %v, want %v", tc.member, got, tc.want)
			}
		})
	}
}

// Remove 只摘掉指定的 member，key 与其余 member 保留。
func TestRedisMemberStore_Remove(t *testing.T) {
	ctx := context.Background()
	client, _ := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	for _, m := range []string{"d1", "d2"} {
		if err := s.Add(ctx, "u1", m, time.Hour); err != nil {
			t.Fatalf("Add(%s): %v", m, err)
		}
	}

	if err := s.Remove(ctx, "u1", "d1"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	got, err := s.Members(ctx, "u1")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if want := []string{"d2"}; !slices.Equal(got, want) {
		t.Errorf("Members = %v, want %v", got, want)
	}
}

// Remove 不带 member 时删除整个 key。
func TestRedisMemberStore_Remove_NoMember_DeletesWholeKey(t *testing.T) {
	ctx := context.Background()
	client, mr := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	for _, m := range []string{"d1", "d2"} {
		if err := s.Add(ctx, "u1", m, time.Hour); err != nil {
			t.Fatalf("Add(%s): %v", m, err)
		}
	}

	if err := s.Remove(ctx, "u1"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if mr.Exists("u1") {
		t.Error("底层 key 仍在，want 被删除")
	}

	got, err := s.Members(ctx, "u1")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Members = %v, want 空", got)
	}
}

// Members 返回前先按 score 清理已过期的 member。
// 注意读路径里带着一次写操作。
func TestRedisMemberStore_Members_DropsExpiredMember(t *testing.T) {
	ctx := context.Background()
	client, mr := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	seedMember(t, mr, "u1", "old", time.Now().Add(-time.Hour))
	seedMember(t, mr, "u1", "fresh", time.Now().Add(time.Hour))

	got, err := s.Members(ctx, "u1")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if want := []string{"fresh"}; !slices.Equal(got, want) {
		t.Errorf("Members = %v, want %v", got, want)
	}
}

// Has 对已过期的 member 返回 false，并顺手把它从 ZSET 摘掉。
func TestRedisMemberStore_Has_ExpiredMember(t *testing.T) {
	ctx := context.Background()
	client, mr := newTestRedis(t)
	s := storex.NewRedisMemberStore(client)

	seedMember(t, mr, "u1", "old", time.Now().Add(-time.Hour))
	seedMember(t, mr, "u1", "fresh", time.Now().Add(time.Hour))

	got, err := s.Has(ctx, "u1", "old")
	if err != nil {
		t.Fatalf("Has: %v", err)
	}
	if got {
		t.Error("过期 member 仍被判为存在")
	}

	members, err := s.Members(ctx, "u1")
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if want := []string{"fresh"}; !slices.Equal(members, want) {
		t.Errorf("Members = %v, want %v", members, want)
	}
}

// ttl <= 0 表示不过期，与 KVStore 一致。
// 走 Add 写入而非 seedMember，是为了覆盖 expiryScore 对非正 TTL 的处理：
// score 落到 +Inf，Members 的过期清理（上界是当前时间戳）碰不到它。
func TestRedisMemberStore_Add_NonPositiveTTL_NeverExpires(t *testing.T) {
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
			client, _ := newTestRedis(t)
			s := storex.NewRedisMemberStore(client)

			if err := s.Add(ctx, "u1", "d1", tc.ttl); err != nil {
				t.Fatalf("Add: %v", err)
			}

			got, err := s.Members(ctx, "u1")
			if err != nil {
				t.Fatalf("Members: %v", err)
			}
			if want := []string{"d1"}; !slices.Equal(got, want) {
				t.Errorf("Members = %v, want %v", got, want)
			}

			has, err := s.Has(ctx, "u1", "d1")
			if err != nil {
				t.Fatalf("Has: %v", err)
			}
			if !has {
				t.Error("非正 TTL 的 member 被判为过期，want 仍有效")
			}
		})
	}
}
