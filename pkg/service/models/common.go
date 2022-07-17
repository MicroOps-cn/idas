package models

type File struct {
	Model
	Name     string `gorm:"type:varchar(128);" `
	Path     string `gorm:"type:varchar(256);" `
	MimiType string `gorm:"type:varchar(50);" `
	Size     int64
}
