package tokenx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// key type 常量
const (
	keyTypeAccess  = "access"
	keyTypeRefresh = "refresh"
	keyTypeSession = "session"
)

type manager struct {
	signer     Signer
	store      Store
	loginMode  LoginMode
	keyPrefix  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// New 创建 Manager 实例。Signer/Store 为必填，nil 时返回 error。
// AccessTokenTTL 零值时兜底 15 分钟，RefreshTokenTTL 零值时兜底 7 天。
func New(cfg Config) (Manager, error) {
	if cfg.Signer == nil {
		return nil, errors.New("tokenx: Signer is required")
	}
	if cfg.Store == nil {
		return nil, errors.New("tokenx: Store is required")
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

	stored, err := m.store.Get(ctx, m.formatKey(keyType, userID, deviceID))
	if err != nil {
		return err
	}
	if stored == "" {
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
	return m.store.Delete(ctx,
		m.formatKey(keyTypeAccess, userID, deviceID),
		m.formatKey(keyTypeRefresh, userID, deviceID),
		m.formatKey(keyTypeSession, userID, deviceID),
	)
}

// RevokeAll 吊销用户全部设备的令牌。先收集所有 key 再一次性删除，避免部分失败。
func (m *manager) RevokeAll(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("tokenx: userID must not be empty")
	}

	var allKeys []string
	for _, typ := range []string{keyTypeAccess, keyTypeRefresh, keyTypeSession} {
		keys, err := m.store.Keys(ctx, m.formatKey(typ, userID, "*"))
		if err != nil {
			return err
		}
		allKeys = append(allKeys, keys...)
	}
	if len(allKeys) > 0 {
		return m.store.Delete(ctx, allKeys...)
	}
	return nil
}

// ListSessions 列出用户所有活跃会话。
func (m *manager) ListSessions(ctx context.Context, userID string) ([]Session, error) {
	if userID == "" {
		return nil, errors.New("tokenx: userID must not be empty")
	}

	keys, err := m.store.Keys(ctx, m.formatKey(keyTypeSession, userID, "*"))
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, 0, len(keys))
	for _, key := range keys {
		data, err := m.store.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		if data == "" {
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
