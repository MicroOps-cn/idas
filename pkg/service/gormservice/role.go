package gormservice

import (
	"context"
	"fmt"

	gogorm "gorm.io/gorm"

	"idas/pkg/global"
	"idas/pkg/service/models"
)

func (c *CommonService) CreateRole(ctx context.Context, role *models.Role) (newRole *models.Role, err error) {
	conn := c.Session(ctx)
	tx := conn.Begin()
	defer tx.Rollback()
	if err = tx.Create(role).Error; err != nil {
		return nil, err
	}
	for _, permission := range role.Permission {
		if err = tx.Model(role).Association("Permission").Append(permission); err != nil {
			return nil, fmt.Errorf("failed to insert permission relationship: %s", err)
		}
	}
	newRole = &models.Role{Model: models.Model{Id: role.Id}}
	err = tx.First(newRole).Error
	if err != nil {
		return nil, err
	} else if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return newRole, nil
}

func (c *CommonService) UpdateRole(ctx context.Context, role *models.Role) (newRole *models.Role, err error) {
	conn := c.Session(ctx)
	tx := conn.Begin()
	defer tx.Rollback()
	if err = tx.Model(role).Association("Permission").Replace(role.Permission); err != nil {
		return nil, err
	}

	newRole = &models.Role{Model: models.Model{Id: role.Id}}
	err = tx.First(newRole).Error
	if err != nil {
		return nil, err
	} else if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return newRole, nil
}

const sqlGetRolesPermission = `
SELECT 
    T0.id, 
    T0.name,  
    T1.role_id
FROM
    t_permission T0
        INNER JOIN
    t_rel_role_permission T1 ON T1.permission_id = T0.id
WHERE
    T1.role_id IN (?)
`

func (c *CommonService) GetRoles(ctx context.Context, keywords string, current, pageSize int64) (count int64, roles []*models.Role, err error) {
	conn := c.Session(ctx)
	tb := conn.Model(&models.Role{})
	if keywords != "" {
		tb = tb.Where("Name LIKE ?", "%"+keywords+"%")
	}
	if err = tb.Count(&count).Error; err != nil {
		return 0, nil, err
	} else if count == 0 {
		return 0, nil, nil
	}
	if err = tb.Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&roles).Error; err != nil {
		return 0, nil, err
	}
	var roleIds []interface{}
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	var permissions []models.RolePermission
	if err = conn.Raw(sqlGetRolesPermission, roleIds).Find(&permissions).Error; err == nil {
		for _, permission := range permissions {
			for idx, role := range roles {
				if role.Id == permission.RoleId {
					pm := permission.Permission
					roles[idx].Permission = append(roles[idx].Permission, &pm)
				}
			}
		}
	}
	return count, roles, nil
}

func (c CommonService) GetPermissions(ctx context.Context, keywords string, current int64, pageSize int64) (count int64, permissions []*models.Permission, err error) {
	conn := c.Session(ctx)
	tb := conn.Model(&models.Permission{})
	if keywords != "" {
		tb = tb.Where("Name LIKE ?", "%"+keywords+"%")
	}
	tb = tb.Where("enable_auth = 1")
	if err = tb.Count(&count).Error; err != nil {
		return 0, nil, err
	} else if count == 0 {
		return 0, nil, nil
	}
	if pageSize == 0 {
		pageSize = 1000
	}
	return count, permissions, tb.Limit(int(pageSize)).Offset(int((current - 1) * pageSize)).Find(&permissions).Error
}

func (c CommonService) DeleteRoles(ctx context.Context, ids []string) error {
	conn := c.Session(ctx)
	tx := conn.Begin()
	defer tx.Rollback()
	if err := conn.Model(models.Role{}).Where("id in (?)", ids).Association("Permission").Delete(); err != nil {
		return err
	} else if err = conn.Model(models.Role{}).Where("id in (?)", ids).Update("is_delete = ?", 1).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

func (c CommonService) CreateOrUpdateRoleByName(ctx context.Context, role *models.Role) error {
	conn := c.Session(ctx)
	if err := conn.Omit("Permission").FirstOrCreate(&role, "name = ?", role.Name).Error; err != nil {
		return err
	}
	return conn.Model(role).Association("Permission").Replace(role.Permission)
}

func (c CommonService) RegisterPermission(ctx context.Context, permissions models.Permissions) error {
	conn := c.Session(ctx)
	for _, p := range permissions {
		fmt.Println(*p)
		var op models.Permission
		if err := conn.Where("name = ?", p.Name).First(&op).Error; err == gogorm.ErrRecordNotFound {
			if err = conn.Create(p).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			p.Id = op.Id
			if op.EnableAuth != p.EnableAuth || op.Description != p.Description || op.ParentId != p.ParentId {
				if err = conn.Select("EnableAuth", "Description", "ParentId").Updates(p).Error; err != nil {
					return err
				}
			}
		}
		for _, child := range p.Children {
			child.ParentId = p.Id
		}
		if err := c.RegisterPermission(context.WithValue(ctx, global.GormConnName, conn), p.Children); err != nil {
			return err
		}
	}
	return nil
}

const sqlAuthorizationByRoleAndMethod = `
SELECT 
    COUNT(1) as count
FROM
    t_role T1
        JOIN
    t_rel_role_permission T2 ON T1.id = T2.role_id
        JOIN
    t_permission T3 ON T3.id = T2.permission_id
WHERE
    T1.name IN ? AND T3.name = ?
`

func (c CommonService) Authorization(ctx context.Context, roles []string, method string) bool {
	var count int64
	if err := c.Session(ctx).Raw(sqlAuthorizationByRoleAndMethod, roles, method).Scan(&count).Error; err != nil {
		return false
	}
	fmt.Println(count)
	return count > 0
}
