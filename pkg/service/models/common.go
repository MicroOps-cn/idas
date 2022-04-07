package models

type File struct {
	Model
	Name     string
	Path     string
	MimiType string
	Size     int64
}
