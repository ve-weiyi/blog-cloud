package permissionx

import "testing"

func TestNormalizeMethod(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"get", "GET"},
		{"  POST  ", "POST"},
		{"", ""},
		{"   ", ""},
		{"*", "*"},
		{"all", "*"},
		{"ALL", "*"},
		{"delete", "DELETE"},
	}
	for _, c := range cases {
		if got := normalizeMethod(c.in); got != c.want {
			t.Errorf("normalizeMethod(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizePath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"/api/v1/users", "/api/v1/users"},
		{"/api/v1/users/", "/api/v1/users"},       // 去尾随斜杠
		{"api/v1/users", "/api/v1/users"},         // 补前导斜杠
		{"/api/v1/users?page=1", "/api/v1/users"}, // 去 query
		{"  /api/v1/users  ", "/api/v1/users"},
		{"/", "/"}, // 根路径保留
		{"", ""},
		{"   ", ""},
		{"/api/v1/users/?a=1", "/api/v1/users"},
	}
	for _, c := range cases {
		if got := normalizePath(c.in); got != c.want {
			t.Errorf("normalizePath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSplitPermID(t *testing.T) {
	cases := []struct {
		in           string
		method, path string
		ok           bool
	}{
		{"GET:/api/v1/users", "GET", "/api/v1/users", true},
		{"POST:/a/b:c", "POST", "/a/b:c", true}, // 只按第一个冒号切分
		{"nocolon", "", "", false},
		{"", "", "", false},
	}
	for _, c := range cases {
		m, p, ok := splitPermID(c.in)
		if m != c.method || p != c.path || ok != c.ok {
			t.Errorf("splitPermID(%q) = (%q,%q,%v), want (%q,%q,%v)",
				c.in, m, p, ok, c.method, c.path, c.ok)
		}
	}
}

func TestPathMatch(t *testing.T) {
	cases := []struct {
		pattern, path string
		want          bool
	}{
		{"/api/v1/users", "/api/v1/users", true},
		{"/api/v1/users", "/api/v1/posts", false},
		{"/api/v1/users/:id", "/api/v1/users/123", true},
		{"/api/v1/users/:id", "/api/v1/users", false},           // :param 不匹配缺失的段
		{"/api/v1/users/:id", "/api/v1/users/123/extra", false}, // 段数必须一致
		{"/api/v1/*", "/api/v1/users/123/extra", true},          // * 吃掉剩余所有段
		{"/api/*", "/api", false},                               // * 不作为零段匹配
		{"/*", "/anything/at/all", true},
		{"/api/v1/users", "/api/v1/users/", true}, // 两侧都按段切分，尾斜杠被 Trim
	}
	for _, c := range cases {
		if got := pathMatch(c.pattern, c.path); got != c.want {
			t.Errorf("pathMatch(%q, %q) = %v, want %v", c.pattern, c.path, got, c.want)
		}
	}
}

func TestMatchPermission(t *testing.T) {
	cases := []struct {
		pattern, target string
		want            bool
	}{
		{"GET:/api/v1/users", "GET:/api/v1/users", true},
		{"GET:/api/v1/users", "POST:/api/v1/users", false}, // 方法必须一致
		{"*:/api/v1/users", "DELETE:/api/v1/users", true},  // pattern 方法为通配
		{"*/api/v1/users", "DELETE:/api/v1/users", false},  // 缺少分隔冒号 -> 拆不出方法
		{"GET:/api/v1/users/:id", "GET:/api/v1/users/42", true},
		{"GET:/api/v1/*", "GET:/api/v1/anything/deep", true},
		{"GET:/a", "GET:/b", false},
		{"GET:/api", "POST:/api/extra", false},
	}
	for _, c := range cases {
		if got := matchPermission(c.pattern, c.target); got != c.want {
			t.Errorf("matchPermission(%q, %q) = %v, want %v", c.pattern, c.target, got, c.want)
		}
	}
}

// apiPermission.Match 是 gorbac 策略侧的入口：p（策略）.Match(req)
func TestApiPermissionMatch(t *testing.T) {
	cases := []struct {
		policy, request string
		want            bool
	}{
		{"GET:/api/v1/users", "GET:/api/v1/users", true},
		{"GET:/api/v1/users/:id", "GET:/api/v1/users/7", true},
		{"*:/api/v1/users", "DELETE:/api/v1/users", true},
		{"GET:/api/v1/users", "POST:/api/v1/users", false},
	}
	for _, c := range cases {
		got := apiPermission(c.policy).Match(apiPermission(c.request))
		if got != c.want {
			t.Errorf("apiPermission(%q).Match(%q) = %v, want %v", c.policy, c.request, got, c.want)
		}
	}
}
