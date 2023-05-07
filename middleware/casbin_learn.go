package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/color"
)

var Enforcer *casbin.Enforcer

func InitCasbinEnforcer() {
	// Your driver and data source
	// 参数true表示若表不存在则自动建表
	a, err1 := gormAdapter.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:3306)/managesys", true)
	Enforcer, _ = casbin.NewEnforcer("configs/model.conf", a)
	if err1 != nil {
		color.Redln(fmt.Sprintf("ERROR!%v", err1.Error()))
	}
	// 加载数据库的权限
	Enforcer.LoadPolicy()
}

func TestCasbin() {
	InitCasbinEnforcer()
	/*
			测试casbin是否能运行

		sub := "alice" // 想要访问资源的用户。
		obj := "data1" // 将被访问的资源。
		act := "read"  // 用户对资源执行的操作。

		ok, err2 := Enforcer.Enforce(sub, obj, act)
		if err2 != nil {
			// 处理err
			color.Redln(fmt.Sprintf("ERROR!%v", err2.Error()))
		}
		if ok == true {
			// 允许alice读取data1
			color.Greenln("Access!")
		} else {
			// 拒绝请求，抛出异常
			color.Redln("Deny!")
		}
	*/

	/*
		对数据库中的policy进行增删改查
	*/

	// 给角色添加权限(policy),若已存在则不会重复添加
	Enforcer.AddPolicy("role_admin", "templates/*", "supreme")
	Enforcer.AddPolicy("role_user", "templates/user", "normal")
	Enforcer.AddPolicy("role_guest", "templates/index.html", "minimal")

	// 根据参数查询信息
	// 第一个参数是匹配起始点的索引，从第二个参数开始按参数顺序查询若参数为""或没填则该位置可为任意值
	info := Enforcer.GetFilteredPolicy(0, "", "templates")
	color.Greenln(info)
	info = Enforcer.GetFilteredPolicy(0, "", "", "supreme")
	color.Greenln(info)
	info = Enforcer.GetFilteredPolicy(0, "role_guest")
	color.Greenln(info)

	/*
		一次性添加多个权限信息
		rules := [][] string {
			[]string {"jack", "data4", "read"},
			[]string {"katy", "data4", "write"},
			[]string {"leyo", "data4", "read"},
			[]string {"ham", "data4", "write"},
		}
		areRulesAdded := Enforcer.AddPolicies(rules)
	*/

	// 根据参数删除信息
	removed, _ := Enforcer.RemovePolicy("alice", "data1", "read")
	color.Greenln(removed)
	/*
		一次性删除多个权限信息
		rules := [][] string {
			[]string {"jack", "data4", "read"},
			[]string {"katy", "data4", "write"},
			[]string {"leyo", "data4", "read"},
			[]string {"ham", "data4", "write"},
		}
		areRulesRemoved := Enforcer.RemovePolicies(rules)
	*/

	// 更新权限信息
	// 首参为旧信息，次参为新信息
	updated, err3 := Enforcer.UpdatePolicy([]string{"nige", "data1", "read"}, []string{"eve", "data1", "write"})
	if err3 != nil {
		color.Redln(fmt.Sprintf("ERROR!%v", err3.Error()))
	}
	color.Greenln(updated)

}

func RBAC() {
	InitCasbinEnforcer()
	// 为用户分配角色(一个用户可以对应多个角色，一个角色可以分给多个用户)
	Enforcer.AddRoleForUser("alice", "role_user")
	// 获取当前用户的角色
	res, err2 := Enforcer.GetRolesForUser("alice")
	if err2 != nil {
		color.Redln(fmt.Sprintf("ERROR!%v", err2.Error()))
	}
	color.Greenln(res)
	// 获取角色的权限
	res1 := Enforcer.GetPermissionsForUser("role_user")
	color.Greenln(res1)

	// 删除指定用户的指定角色
	Enforcer.DeleteRoleForUser("alice", "data1_admin")
}
