package tokenx_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

// mockStore is an in-memory Store for testing Manager logic.
type mockStore struct {
	mu   sync.Mutex
	data map[string]string
}

func newMockStore() *mockStore {
	return &mockStore{data: make(map[string]string)}
}

func (m *mockStore) Set(_ context.Context, key, value string, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *mockStore) Get(_ context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[key], nil
}

func (m *mockStore) Delete(_ context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		delete(m.data, k)
	}
	return nil
}

func (m *mockStore) Keys(_ context.Context, pattern string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	prefix := strings.TrimSuffix(pattern, "*")
	var result []string
	for k := range m.data {
		if strings.HasPrefix(k, prefix) {
			result = append(result, k)
		}
	}
	return result, nil
}

// stubSigner implements Signer minimally for constructor tests.
type stubSigner struct{}

func (s *stubSigner) Sign(_ context.Context, _ tokenx2.TokenClaims) (string, error) { return "", nil }
func (s *stubSigner) Verify(_ context.Context, _ string) (*tokenx2.TokenClaims, error) {
	return nil, nil
}

// stubStore implements Store minimally for constructor tests.
type stubStore struct{}

func (s *stubStore) Set(_ context.Context, _ string, _ string, _ time.Duration) error { return nil }
func (s *stubStore) Get(_ context.Context, _ string) (string, error)                  { return "", nil }
func (s *stubStore) Delete(_ context.Context, _ ...string) error                      { return nil }
func (s *stubStore) Keys(_ context.Context, _ string) ([]string, error)               { return nil, nil }

func TestNew_NilSigner_ReturnsError(t *testing.T) {
	_, err := tokenx2.New(tokenx2.Config{Store: &stubStore{}})
	if err == nil {
		t.Fatal("expected error for nil Signer")
	}
}

func TestNew_NilStore_ReturnsError(t *testing.T) {
	_, err := tokenx2.New(tokenx2.Config{Signer: &stubSigner{}})
	if err == nil {
		t.Fatal("expected error for nil Store")
	}
}

func TestNew_DefaultTTLs(t *testing.T) {
	m, err := tokenx2.New(tokenx2.Config{
		Signer: &stubSigner{},
		Store:  &stubStore{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("expected non-nil Manager")
	}
}

func setupManager(t *testing.T) tokenx2.Manager {
	t.Helper()
	m, err := tokenx2.New(tokenx2.Config{
		Signer: tokenx2.NewHMACSigner("test-secret"),
		Store:  newMockStore(),
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
	if err != tokenx2.ErrTokenEmpty {
		t.Fatalf("expected ErrTokenEmpty, got %v", err)
	}
}

func TestManager_Validate_WrongToken_ReturnsErrTokenInvalid(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	m.Generate(ctx, "user1", "device1")

	err := m.Validate(ctx, "user1", "device1", "wrong-token")
	if err != tokenx2.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestManager_Validate_WrongUser_ReturnsErrTokenExpired(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")

	// HMAC: wrong user → different key → not found → expired
	err := m.Validate(ctx, "user2", "device1", pair.AccessToken)
	if err != tokenx2.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired for wrong user, got %v", err)
	}
}

func TestManager_Validate_AfterRevoke_ReturnsErrTokenExpired(t *testing.T) {
	m := setupManager(t)
	ctx := context.Background()

	pair, _ := m.Generate(ctx, "user1", "device1")
	m.Revoke(ctx, "user1", "device1")

	err := m.Validate(ctx, "user1", "device1", pair.AccessToken)
	if err != tokenx2.ErrTokenExpired {
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
	if err != tokenx2.ErrTokenInvalid {
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
	if err1 != tokenx2.ErrTokenExpired || err2 != tokenx2.ErrTokenExpired {
		t.Fatalf("expected all tokens expired after RevokeAll, got %v / %v", err1, err2)
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

func TestManager_SinglePoint_AutoRevokePrevious(t *testing.T) {
	m, err := tokenx2.New(tokenx2.Config{
		Signer:    tokenx2.NewHMACSigner("test"),
		Store:     newMockStore(),
		LoginMode: tokenx2.SinglePoint,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	pair1, _ := m.Generate(ctx, "user1", "device1")
	pair2, _ := m.Generate(ctx, "user1", "device2")

	// device1 should be revoked (SSO: only latest device active)
	err = m.Validate(ctx, "user1", "device1", pair1.AccessToken)
	if err != tokenx2.ErrTokenExpired {
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
	if err != tokenx2.ErrTokenExpired {
		t.Fatalf("expected token revoked after revoke, got %v", err)
	}
}
