package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/casbin-orm-adapter/v2"
	"gorm.io/gorm"
)

type CasbinEnforcer struct {
	enforcer *casbin.Enforcer
}

func NewCasbinEnforcer(db *gorm.DB, modelText string) (*CasbinEnforcer, error) {
	adapter, err := gormadapter.NewAdapter(db, "casbin_rule", true)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &CasbinEnforcer{enforcer: enforcer}, nil
}

func (e *CasbinEnforcer) Enforce(sub, obj, act string) (bool, error) {
	return e.enforcer.Enforce(sub, obj, act)
}

func (e *CasbinEnforcer) AddPolicy(sub, obj, act string) error {
	_, err := e.enforcer.AddPolicy(sub, obj, act)
	return err
}

func (e *CasbinEnforcer) RemovePolicy(sub, obj, act string) error {
	_, err := e.enforcer.RemovePolicy(sub, obj, act)
	return err
}

func (e *CasbinEnforcer) AddRoleForUser(user, role string) error {
	_, err := e.enforcer.AddGroupingPolicy(user, role)
	return err
}

func (e *CasbinEnforcer) RemoveRoleForUser(user, role string) error {
	_, err := e.enforcer.RemoveGroupingPolicy(user, role)
	return err
}

func (e *CasbinEnforcer) GetRolesForUser(user string) ([]string, error) {
	return e.enforcer.GetRolesForUser(user)
}

func (e *CasbinEnforcer) GetPermissionsForUser(user string) [][]string {
	return e.enforcer.GetPermissionsForUser(user)
}

func (e *CasbinEnforcer) UpdatePolicy(oldSub, oldObj, oldAct, newSub, newObj, newAct string) error {
	_, err := e.enforcer.UpdatePolicy([]string{oldSub, oldObj, oldAct}, []string{newSub, newObj, newAct})
	return err
}

func (e *CasbinEnforcer) ReloadPolicy() error {
	return e.enforcer.LoadPolicy()
}

const DefaultModelText = `
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