package models

type Permission struct {
	Model
	Name     string `gorm:"type:varchar(100);unique" json:"name"`
	Path     string `gorm:"unique" json:"path" `
	ParentId string `gorm:"type:char(32)" json:"parentId" `
}
type Role struct {
	Model
	Name        string        `gorm:"type:varchar(50);unique" json:"name"`
	Description string        `gorm:"type:varchar(250);" json:"describe"`
	Permission  []*Permission `gorm:"association_save_reference:false;association_autoupdate:false;association_autocreate:false;many2many:rel_role_permission" json:"permission,omitempty"`
}

type RolePermission struct {
	Permission
	RoleId string `json:"roleId"`
}
