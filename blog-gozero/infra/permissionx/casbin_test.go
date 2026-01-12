package permissionx

import (
	"fmt"
	"log"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/json-adapter/v2"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
)

// 访问实体 (Subject)，访问资源 (Object) 和访问方法 (Action)
const TestSubjectObjectAction = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`

func TestCasbinRBAC(t *testing.T) {
	// 载入模型
	m, err := model.NewModelFromString(TestSubjectObjectAction)
	if err != nil {
		log.Fatalln("字符串加载模型失败!", err)
	}

	// 载入适配器
	b := []byte{}                         // b stores Casbin policy in JSON bytes.
	adapter := jsonadapter.NewAdapter(&b) // Use b as the data source.

	// 初始化
	enforcer, _ := casbin.NewSyncedCachedEnforcer(m, adapter)
	enforcer.SetExpireTime(60 * 60)

	// 1. 添加用户-角色绑定
	enforcer.AddRoleForUser("alice", "admin") // alice 绑定 admin 角色

	// 2. 角色绑定资源权限
	enforcer.AddPolicy("admin", "data1", "read")  // admin 角色可以读取 data1
	enforcer.AddPolicy("admin", "data1", "write") // admin 角色可以写入 data1

	// 3. 检查权限
	ok, _ := enforcer.Enforce("alice", "data1", "read") // true
	fmt.Println("用户 alice 访问 data1 [read] 是否被允许:", ok)

	ok, _ = enforcer.Enforce("alice", "data1", "write") // true
	fmt.Println("用户 alice 访问 data1 [write] 是否被允许:", ok)

	ok, _ = enforcer.Enforce("alice", "data2", "read") // false
	fmt.Println("用户 alice 访问 data2 [read] 是否被允许:", ok)

	subjects, _ := enforcer.GetAllSubjects()
	t.Log("所有用户 (Subjects):", subjects)

	objects, _ := enforcer.GetAllObjects()
	t.Log("所有资源 (Objects):", objects)

	actions, _ := enforcer.GetAllActions()
	t.Log("所有操作 (Actions):", actions)

	roles, _ := enforcer.GetAllRoles()
	t.Log("所有角色 (Roles):", roles)

	// 5. 查询所有策略
	policies, _ := enforcer.GetPolicy()
	t.Log("所有策略 (Policies):", policies)

	// 6. 查询所有用户-角色绑定
	groupingPolicies, _ := enforcer.GetGroupingPolicy()
	t.Log("所有用户-角色绑定 (Grouping Policies):", groupingPolicies)

	t.Log("json", string(b))
	enforcer.SavePolicy()
	t.Log("json", string(b))
}

// 访问实体 (Subject)，领域(Domain)，访问资源 (Object) 和访问方法 (Action)
const TestSubjectDomainObjectAction = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`

func TestCasbinDomain(t *testing.T) {
	// 载入模型
	m, err := model.NewModelFromString(TestSubjectDomainObjectAction)
	if err != nil {
		log.Fatalln("字符串加载模型失败!", err)
	}

	// 载入适配器
	b := []byte{}                         // b stores Casbin policy in JSON bytes.
	adapter := jsonadapter.NewAdapter(&b) // Use b as the data source.

	// 初始化
	enforcer, _ := casbin.NewSyncedCachedEnforcer(m, adapter)
	enforcer.SetExpireTime(60 * 60)

	// 1. 添加用户-角色绑定
	enforcer.AddRoleForUserInDomain("alice", "admin", "domain1") // alice 绑定 admin 角色

	// 2. 角色绑定资源权限
	enforcer.AddPolicy("admin", "domain1", "data1", "read")  // admin 角色可以读取 data1
	enforcer.AddPolicy("admin", "domain1", "data1", "write") // admin 角色可以写入 data1

	// 3. 检查权限
	ok, _ := enforcer.Enforce("alice", "domain1", "data1", "read") // true
	fmt.Println("用户 alice 访问 data1 [read] 是否被允许:", ok)

	ok, _ = enforcer.Enforce("alice", "domain1", "data1", "write") // true
	fmt.Println("用户 alice 访问 data1 [write] 是否被允许:", ok)

	ok, _ = enforcer.Enforce("alice", "domain1", "data2", "read") // false
	fmt.Println("用户 alice 访问 data2 [read] 是否被允许:", ok)
}

func TestCasbinHolder(t *testing.T) {
	permissionRpc := permissionrpc.NewPermissionRpc(zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{
			"localhost:9999",
		},
	}))
	holder := NewCasbinHolder(
		"localhost:6379",
		permissionRpc,
	)

	err := holder.LoadPolicy()
	t.Log(err)

	err = holder.Enforce("61ef925a-acd9-4209-bda1-8e313aa279c0", "/admin_api/v1/user/get_user_info", "GET")
	t.Log(err)

	enforcer := holder.enforcer
	subjects, _ := enforcer.GetAllSubjects()
	t.Log("所有用户 (Subjects):", subjects)

	objects, _ := enforcer.GetAllObjects()
	t.Log("所有资源 (Objects):", objects)

	actions, _ := enforcer.GetAllActions()
	t.Log("所有操作 (Actions):", actions)

	roles, _ := enforcer.GetAllRoles()
	t.Log("所有角色 (Roles):", roles)

	// 5. 查询所有策略
	policies, _ := enforcer.GetPolicy()
	t.Log("所有策略 (Policies):", policies)

	// 6. 查询所有用户-角色绑定
	groupingPolicies, _ := enforcer.GetGroupingPolicy()
	t.Log("所有用户-角色绑定 (Grouping Policies):", groupingPolicies)
}
