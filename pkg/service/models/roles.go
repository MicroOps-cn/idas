/*
 Copyright © 2022 MicroOps-cn.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package models

import (
	"strings"

	"github.com/MicroOps-cn/idas/config"
)

type Permission struct {
	Model
	Name        string         `gorm:"type:varchar(100);uniqueIndex:idx_t_permission_name" json:"name"`
	Description string         `gorm:"type:varchar(100);" json:"description" `
	ParentId    string         `gorm:"type:char(36)" json:"parentId"`
	EnableAuth  bool           `json:"enableAuth"`
	Children    Permissions    `gorm:"-" json:"children,omitempty"`
	Role        []string       `gorm:"-" json:"role,omitempty"`
	EnableAudit bool           `gorm:"-" json:"enableAudit,omitempty"`
	RateLimit   config.Allower `gorm:"-" json:"rateLimit,omitempty"`
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

func (p Permissions) GetMethod(method string) *Permission {
	for _, permission := range p {
		if permission.Name == method {
			return permission
		}
		if m := permission.Children.GetMethod(method); m != nil {
			return m
		}
	}
	return nil
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
	Name        string        `gorm:"type:varchar(50);uniqueIndex:idx_t_role_name" json:"name"`
	Description string        `gorm:"type:varchar(250);" json:"describe"`
	Type        RoleMeta_Type `json:"type"`
	AppId       string        `json:"app_id"`
	Permission  []*Permission `gorm:"association_save_reference:false;association_autoupdate:false;association_autocreate:false;many2many:rel_role_permission" json:"permission,omitempty"`
}

type RolePermission struct {
	Permission
	RoleId string `json:"roleId"`
}
