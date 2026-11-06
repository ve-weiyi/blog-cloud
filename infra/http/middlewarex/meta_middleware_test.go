package middlewarex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
)

// 平台标识以网关配置为准，客户端声明的值不进入 metadata
func TestCtxMetaMiddleware_AppNameFromGateway(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin-api/v1/ping", nil)
	req.Header.Set(headerx.HeaderAppName, "forged-client")
	req = req.WithContext(metadata.NewOutgoingContext(req.Context(), metadata.MD{}))
	rec := httptest.NewRecorder()

	var got string
	handler := NewCtxMetaMiddleware(headerx.AppNameAdminWeb).Handle(func(_ http.ResponseWriter, r *http.Request) {
		md, _ := metadata.FromOutgoingContext(r.Context())
		if values := md.Get(headerx.HeaderAppName); len(values) > 0 {
			got = values[0]
		}
	})
	handler(rec, req)

	if got != headerx.AppNameAdminWeb {
		t.Fatalf("metadata app name = %q, want %q", got, headerx.AppNameAdminWeb)
	}
}

// 客户端身份头与平台标识都不得透传
func TestCtxMetaMiddleware_DropsClientIdentity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/home", nil)
	req.Header.Set(headerx.HeaderUid, "forged-uid")
	req.Header.Set(headerx.HeaderAuthorization, "Bearer forged-token")
	req.Header.Set(headerx.HeaderAppName, "forged-app")
	req = req.WithContext(metadata.NewOutgoingContext(req.Context(), metadata.MD{}))
	rec := httptest.NewRecorder()

	var uid, authorization, app string
	handler := NewCtxMetaMiddleware(headerx.AppNameBlogWeb).Handle(func(_ http.ResponseWriter, r *http.Request) {
		md, _ := metadata.FromOutgoingContext(r.Context())
		if values := md.Get(headerx.HeaderUid); len(values) > 0 {
			uid = values[0]
		}
		if values := md.Get(headerx.HeaderAuthorization); len(values) > 0 {
			authorization = values[0]
		}
		if values := md.Get(headerx.HeaderAppName); len(values) > 0 {
			app = values[0]
		}
	})
	handler(rec, req)

	if uid != "" {
		t.Fatalf("metadata uid = %q, want empty", uid)
	}
	if authorization != "" {
		t.Fatalf("metadata authorization = %q, want empty", authorization)
	}
	if app != headerx.AppNameBlogWeb {
		t.Fatalf("metadata app name = %q, want %q", app, headerx.AppNameBlogWeb)
	}
}
