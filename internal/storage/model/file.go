package model

import "drivepro/internal/storage"

type File struct {
	ID        string
	Name      string
	Size      int64
	CreatedAt int64
	Type      string
	FolderID  string
	Status    storage.FileStatus
}
