package tokenx_test

import (
	"context"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/storex/tokenx"
)

// mockStore 是 storex.KVStore 的内存实现，用于测试 Manager 逻辑。
type mockStore struct {
	mu        sync.Mutex
	data      map[string]string
	keysCalls int // Keys 被调用的次数，用于断言热路径不做 keyspace 扫描
}

func newMockStore() *mockStore {
	return &mockStore{data: make(map[string]string)}
}

// resetKeysCalls 清空 Keys 调用计数。
func (m *mockStore) resetKeysCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.keysCalls = 0
}

// keysCallCount 返回 Keys 被调用的次数。
func (m *mockStore) keysCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.keysCalls
}

func (m *mockStore) Set(_ context.Context, key, value string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	return nil
}

func (m *mockStore) Get(_ context.Context, key string) (string, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	val, ok := m.data[key]
	return val, ok, nil
}

func (m *mockStore) Delete(_ context.Context, key string, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

// Keys 只支持尾随 "*" 这一种 pattern —— Manager 不用比这更复杂的。
func (m *mockStore) Keys(_ context.Context, pattern string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.keysCalls++

	prefix := strings.TrimSuffix(pattern, "*")

	var keys []string
	for k := range m.data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	return keys, nil
}

// mockMemberStore 是 storex.MemberStore 的内存实现。
// 只保留 Manager 用得上的语义：member 集合 + 按 member 记过期时间。
type mockMemberStore struct {
	mu   sync.Mutex
	data map[string]map[string]time.Time // key -> member -> 过期时间，零值表示不过期
}

func newMockMemberStore() *mockMemberStore {
	return &mockMemberStore{data: make(map[string]map[string]time.Time)}
}

func (m *mockMemberStore) Add(_ context.Context, key, member string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data[key] == nil {
		m.data[key] = make(map[string]time.Time)
	}
	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}
	m.data[key][member] = expireAt

	return nil
}

func (m *mockMemberStore) Members(_ context.Context, key string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	members := make([]string, 0, len(m.data[key]))
	for member, expireAt := range m.data[key] {
		if expireAt.IsZero() || expireAt.After(time.Now()) {
			members = append(members, member)
		}
	}
	// 真实实现按 score 排序，这里排序只为让断言结果确定
	sort.Strings(members)

	return members, nil
}

func (m *mockMemberStore) Has(_ context.Context, key, member string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expireAt, ok := m.data[key][member]
	if !ok {
		return false, nil
	}

	return expireAt.IsZero() || expireAt.After(time.Now()), nil
}

func (m *mockMemberStore) Remove(_ context.Context, key string, members ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(members) == 0 {
		delete(m.data, key)
		return nil
	}
	for _, member := range members {
		delete(m.data[key], member)
	}

	return nil
}

// stubSigner implements Signer minimally for constructor tests.
type stubSigner struct{}

func (s *stubSigner) Sign(_ context.Context, _ tokenx.TokenClaims) (string, error) { return "", nil }
func (s *stubSigner) Verify(_ context.Context, _ string) (*tokenx.TokenClaims, error) {
	return nil, nil
}

// stubStore implements storex.KVStore minimally for constructor tests.
type stubStore struct{}

func (s *stubStore) Set(_ context.Context, _ string, _ string, _ time.Duration) error { return nil }
func (s *stubStore) Get(_ context.Context, _ string) (string, bool, error)            { return "", false, nil }
func (s *stubStore) Delete(_ context.Context, _ string, _ ...string) error            { return nil }
func (s *stubStore) Keys(_ context.Context, _ string) ([]string, error)               { return nil, nil }

// stubMemberStore implements storex.MemberStore minimally for constructor tests.
type stubMemberStore struct{}

func (s *stubMemberStore) Add(_ context.Context, _, _ string, _ time.Duration) error { return nil }
func (s *stubMemberStore) Members(_ context.Context, _ string) ([]string, error)     { return nil, nil }
func (s *stubMemberStore) Has(_ context.Context, _, _ string) (bool, error)          { return false, nil }
func (s *stubMemberStore) Remove(_ context.Context, _ string, _ ...string) error     { return nil }

func TestNew_NilSigner_ReturnsError(t *testing.T) {
	_, err := tokenx.New(tokenx.Config{Store: &stubStore{}, Devices: &stubMemberStore{}})
	if err == nil {
		t.Fatal("expected error for nil Signer")
	}
}

func TestNew_NilStore_ReturnsError(t *testing.T) {
	_, err := tokenx.New(tokenx.Config{Signer: &stubSigner{}, Devices: &stubMemberStore{}})
	if err == nil {
		t.Fatal("expected error for nil Store")
	}
}

// Devices 是必填：没有它 RevokeAll 无从知道该删哪些 key。
// 若为此退回全库扫描，就等于把登录热路径上的 O(全库) 又请了回来。
func TestNew_NilDevices_ReturnsError(t *testing.T) {
	_, err := tokenx.New(tokenx.Config{Signer: &stubSigner{}, Store: &stubStore{}})
	if err == nil {
		t.Fatal("expected error for nil Devices")
	}
}

func TestNew_DefaultTTLs(t *testing.T) {
	m, err := tokenx.New(tokenx.Config{
		Signer:  &stubSigner{},
		Store:   &stubStore{},
		Devices: &stubMemberStore{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
}

// setupManager 造一个多点登录的 Manager，两个存储都用内存 mock。
func setupManager(t *testing.T) tokenx.Manager {
	t.Helper()

	return setupManagerWithMode(t, tokenx.MultiPoint)
}

func setupManagerWithMode(t *testing.T, mode tokenx.LoginMode) tokenx.Manager {
	t.Helper()

	m, err := tokenx.New(tokenx.Config{
		Signer:    tokenx.NewHMACSigner("test-secret"),
		Store:     newMockStore(),
		Devices:   newMockMemberStore(),
		LoginMode: mode,
	})
	if err != nil {
		t.Fatal(err)
	}

	return m
}

func TestManager_Generate_ReturnsTokenPair(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, err := m.Generate(ctx, "user1", "device1")
	if err != nil {
		t.Fatal(err)
	}
	if pair.AccessToken == "" {
		t.Fatal("expected non-empty AccessToken")
	}
	if pair.RefreshToken == "" {
		t.Fatal("expected non-empty RefreshToken")
	}
	if pair.UserID != "user1" {
		t.Fatalf("expected UserID user1, got %s", pair.UserID)
	}
	if pair.DeviceID != "device1" {
		t.Fatalf("expected DeviceID device1, got %s", pair.DeviceID)
	}
	if pair.AccessExpiresIn <= 0 {
		t.Fatal("expected positive AccessExpiresIn")
	}
	if pair.RefreshExpiresIn <= 0 {
		t.Fatal("expected positive RefreshExpiresIn")
	}
	if pair.RefreshExpiresAt.IsZero() {
		t.Fatal("expected non-zero RefreshExpiresAt")
	}
}

func TestManager_Validate_ValidToken_ReturnsNil(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")

	err := m.Validate(ctx, "user1", "device1", pair.AccessToken)
	if err != nil {
		t.Fatalf("expected nil error for valid token, got %v", err)
	}
}

func TestManager_Validate_EmptyToken_ReturnsErrTokenEmpty(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	err := m.Validate(ctx, "user1", "device1", "")
	if err != tokenx.ErrTokenEmpty {
		t.Fatalf("expected ErrTokenEmpty, got %v", err)
	}
}

func TestManager_Validate_WrongToken_ReturnsErrTokenInvalid(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	m.Generate(ctx, "user1", "device1")

	err := m.Validate(ctx, "user1", "device1", "wrong-token")
	if err != tokenx.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestManager_Validate_WrongUser_ReturnsErrTokenExpired(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")

	// HMAC: wrong user → different key → not found → expired
	err := m.Validate(ctx, "user2", "device1", pair.AccessToken)
	if err != tokenx.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired for wrong user, got %v", err)
	}
}

func TestManager_Validate_AfterRevoke_ReturnsErrTokenExpired(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")
	m.Revoke(ctx, "user1", "device1")

	err := m.Validate(ctx, "user1", "device1", pair.AccessToken)
	if err != tokenx.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired after revoke, got %v", err)
	}
}

func TestManager_Refresh_ReturnsNewPair(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")

	newPair, err := m.Refresh(ctx, "user1", "device1", pair.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if newPair.AccessToken == pair.AccessToken {
		t.Fatal("expected new access token after refresh")
	}
	if newPair.RefreshToken == pair.RefreshToken {
		t.Fatal("expected new refresh token after refresh (rotation)")
	}
}

func TestManager_Refresh_OldTokenInvalidated(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")
	m.Refresh(ctx, "user1", "device1", pair.RefreshToken)

	// old refresh token should no longer work (key overwritten by new token)
	_, err := m.Refresh(ctx, "user1", "device1", pair.RefreshToken)
	if err != tokenx.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid for old refresh token, got %v", err)
	}
}

func TestManager_RevokeAll(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair1, _ := m.Generate(ctx, "user1", "device1")
	pair2, _ := m.Generate(ctx, "user1", "device2")

	m.RevokeAll(ctx, "user1")

	// all tokens for user1 should be invalid
	err1 := m.Validate(ctx, "user1", "device1", pair1.AccessToken)
	err2 := m.Validate(ctx, "user1", "device2", pair2.AccessToken)
	if err1 != tokenx.ErrTokenExpired || err2 != tokenx.ErrTokenExpired {
		t.Fatalf("expected all tokens expired after RevokeAll, got %v / %v", err1, err2)
	}
}

// RevokeAll 只碰该用户自己的 key。
// 这是索引取代全库扫描后必须守住的边界：设备清单按 userID 取，天然不越界。
func TestManager_RevokeAll_OnlyTouchesTargetUser(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair1, _ := m.Generate(ctx, "user1", "device1")
	pair2, _ := m.Generate(ctx, "user2", "device1")

	if err := m.RevokeAll(ctx, "user1"); err != nil {
		t.Fatal(err)
	}

	if err := m.Validate(ctx, "user1", "device1", pair1.AccessToken); err != tokenx.ErrTokenExpired {
		t.Fatalf("user1 的令牌未被吊销：%v", err)
	}
	if err := m.Validate(ctx, "user2", "device1", pair2.AccessToken); err != nil {
		t.Fatalf("user2 的令牌被误伤：%v", err)
	}
}

// RevokeAll 与 ListSessions 取设备清单走 devices 索引，不去扫 keyspace。
//
// 这条用例钉的是重构的**理由**而不只是结果：改回 Keys 扫描时，
// 上面那些行为断言全都会照常通过，只有这一条会失败。
func TestManager_RevokeAllAndListSessions_DoNotScanKeyspace(t *testing.T) {
	ctx := context.Background()
	store := newMockStore()

	m, err := tokenx.New(tokenx.Config{
		Signer:  tokenx.NewHMACSigner("test-secret"),
		Store:   store,
		Devices: newMockMemberStore(),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := m.Generate(ctx, "user1", "device1"); err != nil {
		t.Fatal(err)
	}
	store.resetKeysCalls()

	if err := m.RevokeAll(ctx, "user1"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.ListSessions(ctx, "user1"); err != nil {
		t.Fatal(err)
	}

	if n := store.keysCallCount(); n != 0 {
		t.Errorf("RevokeAll/ListSessions 调用了 %d 次 Keys，说明又退回全库扫描了", n)
	}
}

func TestManager_RevokeAll_EmptyUserID(t *testing.T) {
	m := setupManager(t)

	if err := m.RevokeAll(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty userID")
	}
}

func TestManager_ListSessions(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	m.Generate(ctx, "user1", "device1")
	m.Generate(ctx, "user1", "device2")

	sessions, err := m.ListSessions(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestManager_ListSessions_EmptyUserID(t *testing.T) {
	m := setupManager(t)

	if _, err := m.ListSessions(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty userID")
	}
}

// Revoke 之后该设备不该再出现在会话清单里 —— 索引里的 member 也要摘掉。
func TestManager_ListSessions_ClearedAfterRevoke(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	m.Generate(ctx, "user1", "device1")
	m.Generate(ctx, "user1", "device2")

	if err := m.Revoke(ctx, "user1", "device1"); err != nil {
		t.Fatal(err)
	}

	sessions, err := m.ListSessions(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].DeviceID != "device2" {
		t.Errorf("剩余会话的 DeviceID = %q, want %q", sessions[0].DeviceID, "device2")
	}
}

// RevokeAll 之后索引整个清空，否则下次 ListSessions 会去读一堆已不存在的 key。
func TestManager_ListSessions_ClearedAfterRevokeAll(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	m.Generate(ctx, "user1", "device1")
	m.Generate(ctx, "user1", "device2")

	if err := m.RevokeAll(ctx, "user1"); err != nil {
		t.Fatal(err)
	}

	sessions, err := m.ListSessions(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(sessions))
	}
}

func TestManager_SinglePoint_AutoRevokePrevious(t *testing.T) {
	m := setupManagerWithMode(t, tokenx.SinglePoint)
	ctx := context.Background()

	pair1, _ := m.Generate(ctx, "user1", "device1")
	pair2, _ := m.Generate(ctx, "user1", "device2")

	// device1 should be revoked (SSO: only latest device active)
	err := m.Validate(ctx, "user1", "device1", pair1.AccessToken)
	if err != tokenx.ErrTokenExpired {
		t.Fatalf("expected device1 revoked via SSO, got %v", err)
	}

	// device2 should be active
	err = m.Validate(ctx, "user1", "device2", pair2.AccessToken)
	if err != nil {
		t.Fatalf("expected device2 valid, got %v", err)
	}
}

func TestManager_Generate_EmptyUserID(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	_, err := m.Generate(ctx, "", "device1")
	if err == nil {
		t.Fatal("expected error for empty userID")
	}
}

func TestManager_Validate_EmptyUserID(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	err := m.Validate(ctx, "", "device1", "token")
	if err == nil {
		t.Fatal("expected error for empty userID")
	}
}

func TestManager_Revoke_EmptyDeviceID(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	// 空 deviceID 归一化为默认设备，Generate/Validate/Revoke 全链路应正常工作
	pair, err := m.Generate(ctx, "user1", "")
	if err != nil {
		t.Fatalf("expected generate success with empty deviceID, got %v", err)
	}
	if pair.DeviceID != "default" {
		t.Fatalf("expected deviceID normalized to default, got %q", pair.DeviceID)
	}

	err = m.Validate(ctx, "user1", "", pair.AccessToken)
	if err != nil {
		t.Fatalf("expected validate success with empty deviceID, got %v", err)
	}

	err = m.Revoke(ctx, "user1", "")
	if err != nil {
		t.Fatalf("expected revoke success with empty deviceID, got %v", err)
	}

	err = m.Validate(ctx, "user1", "", pair.AccessToken)
	if err != tokenx.ErrTokenExpired {
		t.Fatalf("expected token revoked after revoke, got %v", err)
	}
}
