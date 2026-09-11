package permissionx

import (
	"testing"

	"github.com/mikespook/gorbac/v2"
)

func newTestEnforcer(userRoles map[string][]string, rolePerms map[string][]string) *RbacEnforcer {
	rbac := gorbac.New()
	for role, perms := range rolePerms {
		r := gorbac.NewStdRole(role)
		for _, perm := range perms {
			_ = r.Assign(apiPermission(perm))
		}
		_ = rbac.Add(r)
	}
	return &RbacEnforcer{
		rbac:         rbac,
		userRoles:    userRoles,
		policyLoaded: true,
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
