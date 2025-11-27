package provider

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Casbin Enforcer 是线程安全的
func NewCasbinEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	// 1. 初始化 Gorm 适配器 (它会自动在库里创建 casbin_rule 表)
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 2. 加载模型 (RBAC)
	// 这里使用代码定义的模型，避免依赖文件。
	// r: sub(谁), obj(资源), act(动作)
	// p: 定义策略
	// g: 定义角色组
	// m: 匹配逻辑
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	// 3. 创建 Enforcer
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 4. 加载策略
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	log.Println("✅ Casbin initialized successfully")
	return e, nil
}
