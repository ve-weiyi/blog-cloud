package permissionx

import (
	"context"
	"testing"
	"time"

	"github.com/mikespook/gorbac/v2"
	"google.golang.org/grpc"

	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

// registered 是「已注册接口」全集，用于填充 allPerms。
// Enforce 对未注册的接口按设计放行（见 rbac_enforcer.go 的 unregistered 分支），
// 因此被拒绝的用例必须先把接口登记进来，否则会走放行分支而测不到鉴权。
func newTestEnforcer(userRoles map[string][]string, rolePerms map[string][]string, registered []string) *RbacEnforcer {
	rbac := gorbac.New()
	for role, perms := range rolePerms {
		r := gorbac.NewStdRole(role)
		for _, perm := range perms {
			_ = r.Assign(apiPermission(perm))
		}
		_ = rbac.Add(r)
	}

	allPerms := make(map[string]struct{}, len(registered))
	for _, perm := range registered {
		allPerms[perm] = struct{}{}
	}

	cache := make(map[string]userRoleEntry, len(userRoles))
	for user, roles := range userRoles {
		cache[user] = userRoleEntry{roles: roles, loadedAt: time.Now()}
	}

	return &RbacEnforcer{
		rbac:           rbac,
		userRoles:      cache,
		allPerms:       allPerms,
		userRoleMaxAge: time.Minute, // 测试期间不应过期；零值会让内存缓存立即失效
		policyLoaded:   true,
	}
}

// fakeAccessService 只覆写 LoadPolicy / getUserRoles 用得到的方法；
// 其余靠嵌入接口兜底（一旦被用到即 nil panic，测试会立刻暴露）
type fakeAccessService struct {
	accessservice.AccessService
	roles     []*accessservice.Role
	apis      []*accessservice.Api
	resources map[int64][]int64
}

func (f *fakeAccessService) ListRoles(context.Context, *accessservice.ListRolesRequest, ...grpc.CallOption) (*accessservice.ListRolesResponse, error) {
	return &accessservice.ListRolesResponse{List: f.roles}, nil
}

func (f *fakeAccessService) ListApis(context.Context, *accessservice.ListApisRequest, ...grpc.CallOption) (*accessservice.ListApisResponse, error) {
	return &accessservice.ListApisResponse{List: f.apis}, nil
}

func (f *fakeAccessService) GetRoleResource(_ context.Context, in *accessservice.GetRoleResourceRequest, _ ...grpc.CallOption) (*accessservice.GetRoleResourceResponse, error) {
	return &accessservice.GetRoleResourceResponse{ApiIds: f.resources[in.RoleId]}, nil
}

func (f *fakeAccessService) ListUserRoles(context.Context, *accessservice.ListUserRolesRequest, ...grpc.CallOption) (*accessservice.ListUserRolesResponse, error) {
	return &accessservice.ListUserRolesResponse{}, nil
}

// 禁用一个角色，不应把它独占的接口从「已注册」集合里摘掉。
//
// Enforce 对未注册的接口是直接放行的，所以接口一旦因为角色被禁用而"注销"，
// 就会从「仅该角色可访问」翻转为「对所有人放行」——与禁用角色的预期恰好相反。
func TestLoadPolicyKeepsDisabledRoleApisRegistered(t *testing.T) {
	svc := &fakeAccessService{
		roles: []*accessservice.Role{
			{Id: 1, RoleKey: "admin", Status: enums.RoleStatusNormal},
			{Id: 2, RoleKey: "auditor", Status: enums.RoleStatusDisabled},
		},
		apis: []*accessservice.Api{
			{Id: 10, Method: "GET", Path: "/admin-api/v1/audit"},
			{Id: 11, Method: "GET", Path: "/admin-api/v1/users"},
		},
		resources: map[int64][]int64{
			1: {11},
			2: {10}, // 只分配给已禁用的角色
		},
	}

	e := &RbacEnforcer{pr: svc, rbac: gorbac.New(), userRoles: make(map[string]userRoleEntry)}
	if err := e.LoadPolicy(); err != nil {
		t.Fatalf("LoadPolicy: %v", err)
	}

	if _, ok := e.allPerms["GET:/admin-api/v1/audit"]; !ok {
		t.Fatal("被禁用角色独占的接口从登记集合中消失了——它会变成对所有人放行")
	}

	// 登记了，但不该真的被授权：已禁用的角色不进授权图
	if ok, _ := e.Enforce("someone", "/admin-api/v1/audit", "GET"); ok {
		t.Fatal("被禁用角色不应授予任何权限")
	}
}

func TestEnforce(t *testing.T) {
	e := newTestEnforcer(
		map[string][]string{
			"user1": {"admin"},
			"user2": {"viewer", "editor"},
			"root":  {"root"},
		},
		map[string][]string{
			"admin":  {"GET:/api/v1/users", "GET:/api/v1/users/:id"},
			"viewer": {"GET:/api/v1/posts"},
			"editor": {"POST:/api/v1/posts"},
		},
		[]string{
			"GET:/api/v1/users", "GET:/api/v1/users/:id",
			"GET:/api/v1/posts", "POST:/api/v1/posts",
			// 已注册但未授予任何角色 —— 覆盖"接口存在但无权访问"的拒绝路径
			"POST:/api/v1/users", "DELETE:/api/v1/posts",
		},
	)

	cases := []struct {
		name    string
		user    string
		path    string
		method  string
		wantErr bool
	}{
		{"allowed", "user1", "/api/v1/users", "GET", false},
		{"denied", "user1", "/api/v1/users", "POST", true},
		{"root bypass", "root", "/api/v1/anything", "DELETE", false},
		{"param path", "user1", "/api/v1/users/123", "GET", false},
		{"multiple roles allowed", "user2", "/api/v1/posts", "POST", false},
		{"multiple roles denied", "user2", "/api/v1/posts", "DELETE", true},
		{"empty user", "", "/api/v1/users", "GET", true},
		{"empty path", "user1", "", "GET", true},
		{"empty method", "user1", "/api/v1/users", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := e.Enforce(c.user, c.path, c.method)
			if (err != nil) != c.wantErr {
				t.Errorf("Enforce(%q, %q, %q) error=%v, wantErr=%v", c.user, c.path, c.method, err, c.wantErr)
			}
		})
	}
}
