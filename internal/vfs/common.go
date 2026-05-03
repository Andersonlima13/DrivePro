package vfs

import (
	"drivepro/internal/storage/model"
)

type Directory struct {
	Folders []model.Folder
	Files   []model.File
}
