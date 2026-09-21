package staticfile

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileServerRoutesServe(t *testing.T) {
	dir := t.TempDir()
	// 图片 key 形如 batch_code/station/file（2 级目录）
	os.MkdirAll(filepath.Join(dir, "B-1", "s1"), 0o755)
	if err := os.WriteFile(filepath.Join(dir, "B-1", "s1", "0_0_1.png"), []byte("img-data"), 0o644); err != nil {
		t.Fatal(err)
	}

	routes := FileServerRoutes("/static/images", dir)
	if len(routes) != maxLevel {
		t.Fatalf("routes = %d, want %d", len(routes), maxLevel)
	}
	// 路由逐级加深（go-zero 用 :a/:b/:c 模拟通配）
	wantPaths := []string{
		"/static/images/:a",
		"/static/images/:a/:b",
		"/static/images/:a/:b/:c",
		"/static/images/:a/:b/:c/:d",
		"/static/images/:a/:b/:c/:d/:e",
	}
	for i, r := range routes {
		if r.Method != http.MethodGet {
			t.Errorf("route[%d] method = %s, want GET", i, r.Method)
		}
		if r.Path != wantPaths[i] {
			t.Errorf("route[%d] path = %s, want %s", i, r.Path, wantPaths[i])
		}
	}

	// handler 对任意深度统一 StripPrefix + FileServer：模拟 2 级目录请求
	req := httptest.NewRequest(http.MethodGet, "/static/images/B-1/s1/0_0_1.png", nil)
	rr := httptest.NewRecorder()
	routes[0].Handler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%q)", rr.Code, rr.Body.String())
	}
	if rr.Body.String() != "img-data" {
		t.Fatalf("body = %q, want img-data", rr.Body.String())
	}
}
