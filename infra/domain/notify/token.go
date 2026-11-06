package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 退订令牌的错误分类，调用方据此决定是提示「链接失效」还是「参数非法」。
var (
	ErrTokenMalformed = errors.New("notify: unsubscribe token is malformed")
	ErrTokenInvalid   = errors.New("notify: unsubscribe token signature mismatch")
	ErrTokenExpired   = errors.New("notify: unsubscribe token expired")
	ErrTokenGroup     = errors.New("notify: unsubscribe token group is not allowed")
)

// MinSecretLen 是退订签名密钥的最小长度。
//
// 签名是令牌唯一的防伪手段：空值、示例配置里的占位值（CHANGE_ME）和随手写的短口令
// 都能被猜出，拿到它就能给任意用户签退订令牌。低于此长度一律按「未配置」处理 ——
// 黑名单挡不住下一个人换个占位写法。
const MinSecretLen = 32

// checkSecret 判断密钥是否可用。
func checkSecret(secret string) error {
	if len(secret) < MinSecretLen {
		return fmt.Errorf("notify: unsubscribe secret must be at least %d bytes", MinSecretLen)
	}
	return nil
}

// UnsubscribeClaim 是退订令牌的内容。
type UnsubscribeClaim struct {
	UserID  string
	Group   Group
	Channel Channel
	Expires time.Time
}

// SignUnsubscribe 签发退订令牌。
//
// 令牌自带收件人与分组，签名防伪造，过期时间防重放；强制分组直接拒绝签发 ——
// 邮件正文里根本不该出现它们的退订链接。
func SignUnsubscribe(secret, userID string, group Group, ch Channel, expires time.Time) (string, error) {
	if err := checkSecret(secret); err != nil {
		return "", err
	}
	info, ok := GroupInfoOf(group)
	if !ok {
		return "", fmt.Errorf("notify: unknown group %q", group)
	}
	if info.Mandatory {
		return "", fmt.Errorf("notify: group %q is mandatory", group)
	}
	if userID == "" {
		return "", errors.New("notify: user id is empty")
	}

	payload := strings.Join([]string{
		userID, string(group), string(ch), strconv.FormatInt(expires.Unix(), 10),
	}, "|")
	return encode([]byte(payload)) + "." + encode(sign(secret, payload)), nil
}

// VerifyUnsubscribe 校验令牌并返回它的内容。
func VerifyUnsubscribe(secret, token string, now time.Time) (UnsubscribeClaim, error) {
	var claim UnsubscribeClaim
	if err := checkSecret(secret); err != nil {
		return claim, err
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return claim, ErrTokenMalformed
	}
	payloadBytes, err := decode(parts[0])
	if err != nil {
		return claim, ErrTokenMalformed
	}
	payload := string(payloadBytes)

	sig, err := decode(parts[1])
	if err != nil {
		return claim, ErrTokenMalformed
	}
	if !hmac.Equal(sig, sign(secret, payload)) {
		return claim, ErrTokenInvalid
	}

	fields := strings.Split(payload, "|")
	if len(fields) != 4 {
		return claim, ErrTokenMalformed
	}
	exp, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return claim, ErrTokenMalformed
	}

	claim = UnsubscribeClaim{
		UserID:  fields[0],
		Group:   Group(fields[1]),
		Channel: Channel(fields[2]),
		Expires: time.Unix(exp, 0),
	}
	if claim.UserID == "" {
		return claim, ErrTokenMalformed
	}
	if !claim.Expires.After(now) {
		return claim, ErrTokenExpired
	}

	// 强制分组即便令牌签对了也不能生效：签发端会跳过，这里再拒绝一次
	info, ok := GroupInfoOf(claim.Group)
	if !ok || info.Mandatory {
		return claim, ErrTokenGroup
	}
	return claim, nil
}

func sign(secret, payload string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
