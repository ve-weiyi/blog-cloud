package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/blog-cloud/infra/auth/tokenx"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
)

type AdminAuthMiddleware struct {
	verifier tokenx.Service
}

func NewAdminAuthMiddleware(verifier tokenx.Service) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		verifier: verifier,
	}
}

func (m *AdminAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("AdminAuthMiddleware Handle")

		token := headerx.Token(r.Header, headerx.HeaderAuthorization)
		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, fmt.Sprintf("request header field '%v' is missing", headerx.HeaderAuthorization)))
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
		if tokenx.NormalizeDeviceID(r.Header.Get(headerx.HeaderXDeviceId)) != claims.DeviceID {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "token device mismatch"))
			return
		}

		ctx := metadata.AppendToOutgoingContext(r.Context(),
			headerx.HeaderUid, claims.UserID,
			headerx.HeaderXDeviceId, claims.DeviceID,
			headerx.HeaderToken, token,
		)
		// 同步写入 context value，供直接读取 ctx.Value 的 logic 使用
		ctx = context.WithValue(ctx, headerx.HeaderUid, claims.UserID)
		ctx = context.WithValue(ctx, headerx.HeaderXDeviceId, claims.DeviceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
