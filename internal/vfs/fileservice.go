package vfs

import (
	"drivepro/internal/storage/model"
)

type fileRepository interface {
	ListFiles(folderID string) ([]model.File, error)
}

type FileService struct {
	fileRepo fileRepository
}

func NewFileService(fileRepo fileRepository) *FileService {
	return &FileService{
		fileRepo: fileRepo,
	}
}

func (s *FileService) ListFiles(folderID string) ([]model.File, error) {
	return s.fileRepo.ListFiles(folderID)
}
