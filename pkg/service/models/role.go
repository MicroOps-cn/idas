package models

import "strings"

type Permission struct {
	Model
	Name        string      `gorm:"type:varchar(100);unique" json:"name"`
	Description string      `gorm:"type:varchar(100);" json:"description" `
	ParentId    string      `gorm:"type:char(36)" json:"parentId"`
	EnableAuth  bool        `json:"enableAuth"`
	Children    Permissions `gorm:"-" json:"children,omitempty"`
	Role        []string    `gorm:"-" json:"role,omitempty"`
}

type Permissions []*Permission

func (p Permissions) Get(name string) Permissions {
	var perms Permissions
	for _, permission := range p {
		if permission.Name == name {
			perms = append(perms, permission)
		}
		perms = append(perms, permission.Children.Get(name)...)
	}
	return perms
}

func (p Permissions) GetRoles() Roles {
	var roles Roles
	for _, permission := range p {
		for _, roleName := range permission.Role {
			roleName = strings.TrimSpace(roleName)
			if role := roles.Get(strings.TrimSpace(roleName)); role != nil {
				role.Permission = append(role.Permission, permission)
			} else {
				roles = append(roles, &Role{Name: roleName, Permission: Permissions{permission}})
			}
		}
		childRoles := permission.Children.GetRoles()
		for _, childRole := range childRoles {
			if role := roles.Get(childRole.Name); role != nil {
				role.Permission = append(role.Permission, childRole.Permission...)
			} else {
				roles = append(roles, &Role{Name: childRole.Name, Permission: childRole.Permission})
			}
		}
	}
	return roles
}

type Roles []*Role

func (r Roles) Get(name string) *Role {
	for _, role := range r {
		if role.Name == name {
			return role
		}
	}
	return nil
}

type Role struct {
	Model
	Name        string        `gorm:"type:varchar(50);unique" json:"name"`
	Description string        `gorm:"type:varchar(250);" json:"describe"`
	Type        RoleMeta_Type `json:"type"`
	Permission  []*Permission `gorm:"association_save_reference:false;association_autoupdate:false;association_autocreate:false;many2many:rel_role_permission" json:"permission,omitempty"`
}

type RolePermission struct {
	Permission
	RoleId string `json:"roleId"`
}
