package model

type Folder struct {
	ID        string
	Name      string
	ParentID  *string
	CreatedAt int64
}
