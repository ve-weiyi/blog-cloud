package rbac

import (
	"log"
	"path"
	"reflect"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/initest"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

var enforcer *casbin.SyncedEnforcer
var db *gorm.DB

func Init() {
	initest.Init(path.Join(global.GetRuntimeRoot() + "server/config.yaml"))
	db = global.DB
	adapter, err := gormadapter.NewAdapterByDB(db)
	m, err := model.NewModelFromString(SubjectDomainObjectAction)
	if err != nil {
		zap.L().Error("字符串加载模型失败!", zap.Error(err))
		return
	}
	//会自动生成数据库表
	enforcer, err = casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		log.Println("err-->", err)
		return
	}
	//enforcer.SetExpireTime(60 * 60 * 100)
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Println("err-->", err)
		return
	}

}

func TestCasbin(t *testing.T) {
	//清理数据库

	//角色、访问api、请求方法

	//policy, err := enforcer.AddPolicy("admin", "blog", "home", "login")
	//log.Println("-->", policy, err)
	//
	//enforcer.AddPolicy("admin", "soundcore", "home", "login")
	//enforcer.AddPolicy("admin", "eufy", "home", "login")
	//enforcer.AddPolicy("admin", "appliance", "home", "login")
	//enforcer.AddPolicy("super-admin", "soundcore", "home", "login")

	//enforcer.AddRoleForUserInDomain("alice", "admin", "eufy")

	//force, err := enforcer.Enforce("admin", "blog", "home", "login")
	//log.Println("-->", force)

	//为角色添加规则
	enforcer.AddPolicy("admin", "blog", "/api/v1/roles", "GET")
	enforcer.AddPolicy("admin", "blog", "/api/v1/role/list", "GET")
	enforcer.AddPolicy("admin", "blog", "/api/v1/role/create", "POST")
	enforcer.AddPolicy("admin", "blog", "/api/v1/role/update", "PUT")
	enforcer.AddPolicy("admin", "blog", "/api/v1/role/delete", "DELETE")
	enforcer.AddPolicy("admin", "blog", "/api/v1/role/update_menus", "POST")
	enforcer.AddPolicy("admin", "blog", "/api/v1/role/update_resources", "POST")

	enforcer.AddPolicy("admin", "blog", "/api/v1/user/list", "GET")
	enforcer.AddPolicy("admin", "blog", "/api/v1/user/update_roles", "POST")
	log.Println(enforcer.Enforce("admin", "blog", "/api/v1/user/update_roles", "POST"))

	//为用户添加角色
	//enforcer.AddRoleForUserInDomain("ve", "admin", "blog")
	//log.Println(enforcer.Enforce("ve", "blog", "/api/v1/role/list", "GET"))

	//为用户添加单独的权限
	//enforcer.AddPermissionForUser("hi", "blog", "/api/v1/role/list", "GET")
	//log.Println(enforcer.Enforce("hi", "blog", "/api/v1/role/list", "GET"))

	log.Println("-->", enforcer.GetPolicy())
	log.Println("-->", enforcer.GetAllSubjects())
	log.Println("-->", enforcer.GetAllRoles())
	log.Println(enforcer.GetAllDomains())
	log.Println(enforcer.GetAllActions())
	log.Println(enforcer.SavePolicy())
	//用户添加角色
	//enforcer.AddRoleForUser("zhangsan", "member") //这是单条添加用户
}

func TestAddAllPolicy(t *testing.T) {
	//policy, err := enforcer.AddPolicy("admin", "blog", "home", "login")
	//清理数据库
	enforcer.ClearPolicy()
	enforcer.SavePolicy()

	ctx := svc.NewRepositoryContext(&global.CONFIG)
	rp := repository.NewRepository(ctx)

	re := NewCachedEnforcer(global.DB)
	policy, err := re.DeleteRolePolicy("admin", "blog")
	if err != nil {
		return
	}
	data, err := rp.RoleRepository.FindRoleApis(nil, 1)
	rolePolicy, err := re.AddRolePolicy("admin", "blog", data)
	if err != nil {
		return
	}

	log.Println(policy, rolePolicy)
	//ResetAllPolicy(db, enforcer)
}

func ResetAllPolicy(db *gorm.DB, rbac *casbin.SyncedEnforcer) {
	//角色、域、页面、api、method
	type name struct {
		RoleId string
		Domain string
		Page   string
		Api    string
		reflect.Method
	}

	//角色
	var roles []entity.Role
	err := db.Find(&roles).Error
	if err != nil {
		return
	}

	log.Println("--->", jsonconv.ObjectToJsonIndent(roles))
	//域
	for _, role := range roles {
		menus, err := ResetRoleMenuPolicy(db, role.ID)
		if err != nil {
			log.Println("err-->", err)
			return
		}

		log.Println("--->", jsonconv.ObjectToJsonIndent(menus))
		for _, item := range menus {
			rbac.AddPolicy(role.RoleComment, role.RoleDomain, item.Name, item.Path)
		}

	}

	log.Println("--->", rbac.GetPolicy())
}

// 获取角色菜单权限
func ResetRoleMenuPolicy(db *gorm.DB, roleId int) ([]*response.UserMenuDTO, error) {
	var urs []entity.RoleMenu
	err := db.Where("role_id = ?", roleId).Find(&urs).Error
	if err != nil {
		return nil, err
	}

	var mids []int
	for _, item := range urs {
		mids = append(mids, item.MenuID)
	}

	var menus []entity.Menu
	err = db.Where("id in (?)", mids).Find(&menus).Error
	if err != nil {
		return nil, err
	}

	var res []*response.UserMenuDTO
	for _, item := range menus {
		menu := &response.UserMenuDTO{
			Id:        item.ID,
			Name:      item.Name,
			Path:      item.Path,
			Component: item.Component,
			Icon:      item.Icon,
			IsHidden:  item.IsHidden,
		}
		res = append(res, menu)
	}

	return res, nil
}
