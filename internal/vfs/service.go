package vfs

import (
	"drivepro/internal/storage/repository"
)

type Service struct {
	fileRepo   *repository.FileRepository
	folderRepo *repository.FolderRepository
}

func NewService(
	fileRepo *repository.FileRepository,
	folderRepo *repository.FolderRepository,
) *Service {
	return &Service{
		fileRepo:   fileRepo,
		folderRepo: folderRepo,
	}
}
