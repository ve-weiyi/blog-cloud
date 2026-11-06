package presencex

import (
	"context"
	"sort"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/store/storex"
	"github.com/ve-weiyi/stompws/server"
)

// Store 是在线状态的领域存储，实现 stompws 的 server.PresenceManager：
// hub 只负责报告连接的来去，不关心这份状态存在哪。
//
// 一个 service 一个键 —— 键归谁，只有谁往里 Add / Remove。只读别人名册的用法
// （如后台读站点的在线名单）只调 List / PresenceStats：它不改「谁在线」这个语义
// （不增不删有效成员），但会顺手回收已过期的成员 —— 那是 storex 的既有行为，
// 也是这套存储不需要清扫协程的原因。
type Store struct {
	members storex.MemberStore
	key     string
	ttl     time.Duration
}

// NewStore 建一个在线状态存储。ttl <= 0 时取 DefaultTTL。
func NewStore(members storex.MemberStore, key string, ttl time.Duration) *Store {
	if ttl <= 0 {
		ttl = DefaultTTL
	}

	return &Store{members: members, key: key, ttl: ttl}
}

// TTL 给装配处配 hub 的续期存活时长用 —— 存活时长只此一处定义，
// 存储与 hub 不会各写一个数。
func (s *Store) TTL() time.Duration { return s.ttl }

// Add 登记一条连接并续期；同一 (login, ref) 重复调用是幂等的，续期走的就是这条路。
func (s *Store) Add(ctx context.Context, login, ref string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = s.ttl
	}

	return s.members.Add(ctx, s.key, EncodeMember(login, ref), ttl)
}

// Remove 摘除一条连接；该 login 还有其他 ref 时，用户仍在名单里。
func (s *Store) Remove(ctx context.Context, login, ref string) error {
	return s.members.Remove(ctx, s.key, EncodeMember(login, ref))
}

func (s *Store) PresenceStats(ctx context.Context) (server.PresenceStats, error) {
	users, err := s.List(ctx)
	if err != nil {
		return server.PresenceStats{}, err
	}

	stats := server.PresenceStats{NumUsers: len(users)}
	for _, user := range users {
		stats.NumClients += user.Refs
	}

	return stats, nil
}

// Presence 列出在线的连接 —— 一个用户开着几条连接就有几条。
func (s *Store) Presence(ctx context.Context) ([]server.PresenceMember, error) {
	expiries, err := s.members.Expiries(ctx, s.key)
	if err != nil {
		return nil, err
	}

	members := make([]server.PresenceMember, 0, len(expiries))
	for member, expireAt := range expiries {
		login, ref, ok := DecodeMember(member)
		if !ok {
			continue
		}
		members = append(members, server.PresenceMember{
			Login:        login,
			Ref:          ref,
			LastActiveAt: expireToLastActive(expireAt, s.ttl),
		})
	}
	sort.Slice(members, func(i, j int) bool {
		if members[i].Login != members[j].Login {
			return members[i].Login < members[j].Login
		}
		return members[i].Ref < members[j].Ref
	})

	return members, nil
}

// List 按用户列出在线名单，供业务侧展示。
func (s *Store) List(ctx context.Context) ([]OnlineUser, error) {
	expiries, err := s.members.Expiries(ctx, s.key)
	if err != nil {
		return nil, err
	}

	return aggregate(expiries, s.ttl), nil
}

// Close 无自身资源要释放：过期成员由 storex 在读路径上惰性清理，这里不起后台协程。
func (s *Store) Close() error { return nil }

var _ server.PresenceManager = (*Store)(nil)
