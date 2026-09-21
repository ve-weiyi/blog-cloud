package tokenx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

// key type 常量
const (
	keyTypeAccess  = "access"
	keyTypeRefresh = "refresh"
	keyTypeSession = "session"
	keyTypeDevices = "devices" // 设备索引，按 userID 聚合，不带 deviceID
)

type manager struct {
	signer     Signer
	store      storex.KVStore
	devices    storex.MemberStore
	loginMode  LoginMode
	keyPrefix  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// New 创建 Manager 实例。Signer/Store/Devices 为必填，nil 时返回 error。
// AccessTokenTTL 零值时兜底 15 分钟，RefreshTokenTTL 零值时兜底 7 天。
func New(cfg Config) (Manager, error) {
	if cfg.Signer == nil {
		return nil, errors.New("tokenx: Signer is required")
	}
	if cfg.Store == nil {
		return nil, errors.New("tokenx: Store is required")
	}
	if cfg.Devices == nil {
		return nil, errors.New("tokenx: Devices is required")
	}
	accessTTL := cfg.AccessTokenTTL
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	refreshTTL := cfg.RefreshTokenTTL
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &manager{
		signer:     cfg.Signer,
		store:      cfg.Store,
		devices:    cfg.Devices,
		loginMode:  cfg.LoginMode,
		keyPrefix:  cfg.KeyPrefix,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

// formatKey 构造 Store key：prefix:type:userID:deviceID
func (m *manager) formatKey(keyType, userID, deviceID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", m.keyPrefix, keyType, userID, deviceID)
}

// formatDeviceKey 构造设备索引 key：prefix:devices:userID
// 索引按 userID 聚合，因此不带 deviceID —— deviceID 是它的 member。
func (m *manager) formatDeviceKey(userID string) string {
	return fmt.Sprintf("%s:%s:%s", m.keyPrefix, keyTypeDevices, userID)
}

// deviceKeys 返回某台设备名下的全部令牌 key。
func (m *manager) deviceKeys(userID, deviceID string) []string {
	return []string{
		m.formatKey(keyTypeAccess, userID, deviceID),
		m.formatKey(keyTypeRefresh, userID, deviceID),
		m.formatKey(keyTypeSession, userID, deviceID),
	}
}

// deleteKeys 一次删除一组 key，空集合直接跳过。
// storex 的 Delete 首参必填（不允许一个 key 都不传），这里做那层转发。
func (m *manager) deleteKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return m.store.Delete(ctx, keys[0], keys[1:]...)
}

// defaultDeviceID 空 deviceID 归一化后的默认设备标识。
const defaultDeviceID = "default"

// normalizeDeviceID 空 deviceID 归一化为默认设备，保证 store key 与 claims 比对一致。
func (m *manager) normalizeDeviceID(deviceID string) string {
	if deviceID == "" {
		return defaultDeviceID
	}
	return deviceID
}

// validateIdentity 校验 userID 非空，deviceID 允许为空（空值由 normalizeDeviceID 归一化处理）。
func (m *manager) validateIdentity(userID, deviceID string) error {
	if userID == "" {
		return errors.New("tokenx: userID must not be empty")
	}
	return nil
}

// verifyToken 通用令牌校验。校验通过返回 nil。
func (m *manager) verifyToken(ctx context.Context, userID, deviceID, token, keyType string) error {
	if token == "" {
		return ErrTokenEmpty
	}

	claims, err := m.signer.Verify(ctx, token)
	if err != nil {
		return err
	}
	if claims != nil {
		if claims.UserID != userID || claims.DeviceID != deviceID {
			return ErrTokenInvalid
		}
		if claims.TokenType != "" && claims.TokenType != keyType {
			return ErrTokenInvalid
		}
	}

	stored, found, err := m.store.Get(ctx, m.formatKey(keyType, userID, deviceID))
	if err != nil {
		return err
	}
	if !found {
		return ErrTokenExpired
	}
	if stored != token {
		return ErrTokenInvalid
	}
	return nil
}

// Generate 生成令牌对。
func (m *manager) Generate(ctx context.Context, userID, deviceID string) (*TokenPair, error) {
	deviceID = m.normalizeDeviceID(deviceID)
	if err := m.validateIdentity(userID, deviceID); err != nil {
		return nil, err
	}

	now := time.Now()

	accessClaims := TokenClaims{
		UserID:    userID,
		DeviceID:  deviceID,
		TokenType: keyTypeAccess,
		IssuedAt:  now,
		ExpiresAt: now.Add(m.accessTTL),
	}
	refreshClaims := TokenClaims{
		UserID:    userID,
		DeviceID:  deviceID,
		TokenType: keyTypeRefresh,
		IssuedAt:  now,
		ExpiresAt: now.Add(m.refreshTTL),
	}

	accessToken, err := m.signer.Sign(ctx, accessClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := m.signer.Sign(ctx, refreshClaims)
	if err != nil {
		return nil, err
	}

	// SSO：吊销用户所有旧令牌
	if m.loginMode == SinglePoint {
		if err := m.RevokeAll(ctx, userID); err != nil {
			return nil, err
		}
	}

	accessKey := m.formatKey(keyTypeAccess, userID, deviceID)
	refreshKey := m.formatKey(keyTypeRefresh, userID, deviceID)
	sessionKey := m.formatKey(keyTypeSession, userID, deviceID)

	if err := m.store.Set(ctx, accessKey, accessToken, m.accessTTL); err != nil {
		return nil, err
	}
	if err := m.store.Set(ctx, refreshKey, refreshToken, m.refreshTTL); err != nil {
		m.store.Delete(ctx, accessKey) // 回滚
		return nil, err
	}

	session := Session{
		UserID:    userID,
		DeviceID:  deviceID,
		CreatedAt: now,
		ExpiresAt: now.Add(m.refreshTTL),
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		m.store.Delete(ctx, accessKey, refreshKey) // 回滚
		return nil, fmt.Errorf("tokenx: marshal session: %w", err)
	}
	if err := m.store.Set(ctx, sessionKey, string(sessionJSON), m.refreshTTL); err != nil {
		m.store.Delete(ctx, accessKey, refreshKey) // 回滚
		return nil, err
	}

	// 设备索引写在最后：失败时的回滚就是上面那三个 key，不新增回滚分支。
	// 索引缺项会让 RevokeAll 漏掉该设备（危险方向），所以这里必须失败即回滚。
	if err := m.devices.Add(ctx, m.formatDeviceKey(userID), deviceID, m.refreshTTL); err != nil {
		m.store.Delete(ctx, accessKey, refreshKey, sessionKey) // 回滚
		return nil, err
	}

	return &TokenPair{
		UserID:           userID,
		DeviceID:         deviceID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  now.Add(m.accessTTL),
		AccessExpiresIn:  int64(m.accessTTL.Seconds()),
		RefreshExpiresAt: now.Add(m.refreshTTL),
		RefreshExpiresIn: int64(m.refreshTTL.Seconds()),
	}, nil
}

// Validate 校验 AccessToken。
func (m *manager) Validate(ctx context.Context, userID, deviceID, accessToken string) error {
	deviceID = m.normalizeDeviceID(deviceID)
	if err := m.validateIdentity(userID, deviceID); err != nil {
		return err
	}
	return m.verifyToken(ctx, userID, deviceID, accessToken, keyTypeAccess)
}

// Refresh 校验 RefreshToken 并轮换。
// Generate 会覆盖 refresh key，旧 RefreshToken 自动失效。
func (m *manager) Refresh(ctx context.Context, userID, deviceID, refreshToken string) (*TokenPair, error) {
	deviceID = m.normalizeDeviceID(deviceID)
	if err := m.validateIdentity(userID, deviceID); err != nil {
		return nil, err
	}
	if err := m.verifyToken(ctx, userID, deviceID, refreshToken, keyTypeRefresh); err != nil {
		return nil, err
	}
	return m.Generate(ctx, userID, deviceID)
}

// Revoke 吊销指定设备的所有令牌。
func (m *manager) Revoke(ctx context.Context, userID, deviceID string) error {
	deviceID = m.normalizeDeviceID(deviceID)
	if err := m.validateIdentity(userID, deviceID); err != nil {
		return err
	}

	if err := m.deleteKeys(ctx, m.deviceKeys(userID, deviceID)); err != nil {
		return err
	}

	// 摘掉索引里的 deviceID，否则 RevokeAll 会去删已经不存在的 key
	return m.devices.Remove(ctx, m.formatDeviceKey(userID), deviceID)
}

// RevokeAll 吊销用户全部设备的令牌。
// 设备清单取自 devices 索引，因此只访问该用户自己的 key，不做全库扫描。
func (m *manager) RevokeAll(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("tokenx: userID must not be empty")
	}

	deviceKey := m.formatDeviceKey(userID)
	deviceIDs, err := m.devices.Members(ctx, deviceKey)
	if err != nil {
		return err
	}
	if len(deviceIDs) == 0 {
		return nil
	}

	keys := make([]string, 0, len(deviceIDs)*3)
	for _, deviceID := range deviceIDs {
		keys = append(keys, m.deviceKeys(userID, deviceID)...)
	}

	// 先删令牌再清索引：反过来的话，令牌删除失败而索引已清空，
	// 这批令牌就再也没人找得到，只能等 TTL 自然过期。
	if err := m.deleteKeys(ctx, keys); err != nil {
		return err
	}

	return m.devices.Remove(ctx, deviceKey)
}

// ListSessions 列出用户所有活跃会话。
func (m *manager) ListSessions(ctx context.Context, userID string) ([]Session, error) {
	if userID == "" {
		return nil, errors.New("tokenx: userID must not be empty")
	}

	deviceIDs, err := m.devices.Members(ctx, m.formatDeviceKey(userID))
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(deviceIDs))
	for _, deviceID := range deviceIDs {
		data, found, err := m.store.Get(ctx, m.formatKey(keyTypeSession, userID, deviceID))
		if err != nil {
			return nil, err
		}
		if !found {
			// 会话已先于索引过期，跳过
			continue
		}

		var s Session
		if err := json.Unmarshal([]byte(data), &s); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}
