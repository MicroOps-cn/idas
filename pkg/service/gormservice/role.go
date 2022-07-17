package gormservice

import (
	"context"
	"fmt"
	"idas/pkg/service/models"
)

func (c *CommonService) CreateRole(ctx context.Context, role *models.Role) (err error) {
	conn := c.Session(ctx)
	tx := conn.Begin()
	defer tx.Rollback()
	if err = tx.Create(role).Error; err != nil {
		return err
	}
	for _, permission := range role.Permission {
		if err = tx.Model(role).Association("Permission").Append(permission); err != nil {
			return fmt.Errorf("failed to insert permission relationship: %s", err)
		}
	}
	return tx.Commit().Error
}

func (c *CommonService) UpdateRole(ctx context.Context, role *models.Role) (err error) {
	conn := c.Session(ctx)
	tx := conn.Begin()
	defer tx.Rollback()
	if err = tx.Model(role).Association("Permission").Replace(role.Permission); err != nil {
		return err
	}

	return tx.Commit().Error
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

func (c *CommonService) GetRoles(ctx context.Context, keywords string, current int, pageSize int) (count int64, roles []*models.Role, err error) {
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
	if err = tb.Limit(pageSize).Offset((current - 1) * pageSize).Find(&roles).Error; err != nil {
		return 0, nil, err
	}
	var roleIds []interface{}
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}
	var permissions []models.RolePermission
	if err := conn.Raw(sqlGetRolesPermission, roleIds).Find(&permissions).Error; err == nil {
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
func (c CommonService) GetPermissions(ctx context.Context, keywords string, current int, pageSize int) (count int64, permissions []*models.Permission, err error) {
	conn := c.Session(ctx)
	tb := conn.Model(&models.Permission{})
	if keywords != "" {
		tb = tb.Where("Name LIKE ?", "%"+keywords+"%")
	}
	if err = tb.Count(&count).Error; err != nil {
		return 0, nil, err
	} else if count == 0 {
		return 0, nil, nil
	}
	if pageSize == 0 {
		pageSize = 1000
	}
	return count, permissions, tb.Limit(pageSize).Offset((current - 1) * pageSize).Find(&permissions).Error
}

func (c CommonService) RemoveRoles(ctx context.Context, ids []string) error {
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
