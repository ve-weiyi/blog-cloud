package middleware

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
	"github.com/ve-weiyi/vkit/x/cryptox"
)

// guestTokenWindow 游客请求时间戳允许的偏差，防止旧请求被重放。
const guestTokenWindow = 10 * time.Minute

// OptionalUserAuthMiddleware 允许登录用户或带签名的游客通过。
//
// 有令牌时按登录用户校验并写入身份；没有令牌时要求合法的设备签名。
// AI 助手对所有访客开放，因此不能直接复用强制的 UserAuth。
type OptionalUserAuthMiddleware struct {
	verifier tokenx.Service
}

func NewOptionalUserAuthMiddleware(verifier tokenx.Service) *OptionalUserAuthMiddleware {
	return &OptionalUserAuthMiddleware{verifier: verifier}
}

func (m *OptionalUserAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("OptionalUserAuthMiddleware Handle")

		deviceID := requestDeviceID(r)
		token := headerx.Token(r.Header, headerx.HeaderAuthorization)

		if token == "" {
			if err := validateGuestDevice(r, deviceID); err != nil {
				responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, err.Error()))
				return
			}
			ctx := metadata.AppendToOutgoingContext(r.Context(), headerx.HeaderXDeviceId, deviceID)
			ctx = context.WithValue(ctx, headerx.HeaderXDeviceId, deviceID)
			ctx = aictx.WithIdentity(ctx, aictx.Identity{DeviceId: deviceID, Guest: true})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, err := m.verifier.Validate(r.Context(), token)
		if err != nil {
			if errors.Is(err, tokenx.ErrTokenExpired) {
				responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeLoginExpired, err.Error()))
				return
			}
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, err.Error()))
			return
		}

		// 令牌绑定设备：空设备标识按 default 归一化后也必须与令牌一致。
		if tokenx.NormalizeDeviceID(requestDeviceID(r)) != claims.DeviceID {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "token device mismatch"))
			return
		}

		ctx := metadata.AppendToOutgoingContext(r.Context(),
			headerx.HeaderUid, claims.UserID,
			headerx.HeaderXDeviceId, claims.DeviceID,
			headerx.HeaderToken, token,
		)
		ctx = context.WithValue(ctx, headerx.HeaderUid, claims.UserID)
		ctx = context.WithValue(ctx, headerx.HeaderXDeviceId, claims.DeviceID)
		ctx = aictx.WithIdentity(ctx, aictx.Identity{UserId: claims.UserID, DeviceId: claims.DeviceID})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// validateGuestDevice 校验游客设备标识与签名，并要求时间戳在允许窗口内。
func validateGuestDevice(r *http.Request, deviceID string) error {
	if deviceID == "" {
		return fmt.Errorf("缺少登录令牌或设备标识")
	}

	headerDevice := r.Header.Get(headerx.HeaderXDeviceId)
	timestamp := r.Header.Get(headerx.HeaderTimestamp)
	token := r.Header.Get(headerx.HeaderXDeviceToken)
	if headerDevice == "" || timestamp == "" || token == "" {
		return fmt.Errorf("游客请求缺少设备签名")
	}
	if headerDevice != deviceID {
		return fmt.Errorf("设备标识不一致")
	}

	sec, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("设备签名时间戳无效")
	}
	if diff := time.Since(time.Unix(sec, 0)); diff > guestTokenWindow || diff < -guestTokenWindow {
		return fmt.Errorf("设备签名已过期")
	}

	expected := cryptox.Sha256v(headerDevice, timestamp)
	if subtle.ConstantTimeCompare([]byte(expected), []byte(token)) != 1 {
		return fmt.Errorf("无效请求,游客签名错误")
	}
	return nil
}
