package vfs

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"drivepro/internal/storage/model"
)

var errFolderNotFound = errors.New("folder not found")

type Directory struct {
	Folders []model.Folder
	Files   []model.File
}

type fileRepository interface {
	ListFiles(folderID string) ([]model.File, error)
}

type folderRepository interface {
	CreateFolder(name string, parentID *string) (*model.Folder, error)
	ListFolders(parentID *string) ([]model.Folder, error)
}

type Service struct {
	fileRepo   fileRepository
	folderRepo folderRepository
}

func NewService(
	fileRepo fileRepository,
	folderRepo folderRepository,
) *Service {
	return &Service{
		fileRepo:   fileRepo,
		folderRepo: folderRepo,
	}
}

func (s *Service) CreateFolderByPath(folderPath string) error {
	parts, err := splitFolderPath(folderPath)
	if err != nil {
		return err
	}

	if len(parts) == 0 {
		return errors.New("cannot create root folder")
	}

	var parentID *string

	for i, name := range parts {
		folder, err := s.findChildFolder(parentID, name)
		if err == nil {
			if i == len(parts)-1 {
				return fmt.Errorf("folder already exists: %s", folderPath)
			}

			parentID = &folder.ID
			continue
		}

		if !errors.Is(err, errFolderNotFound) {
			return err
		}

		folder, err = s.folderRepo.CreateFolder(name, parentID)
		if err != nil {
			return err
		}

		parentID = &folder.ID
	}

	return nil
}

func (s *Service) GetFolderByPath(folderPath string) (*model.Folder, error) {
	parts, err := splitFolderPath(folderPath)
	if err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return rootFolder(), nil
	}

	var parentID *string
	var folder *model.Folder

	for _, name := range parts {
		folder, err = s.findChildFolder(parentID, name)
		if err != nil {
			if errors.Is(err, errFolderNotFound) {
				return nil, fmt.Errorf("folder not found: %s", folderPath)
			}

			return nil, err
		}

		parentID = &folder.ID
	}

	return folder, nil
}

func (s *Service) ListDirectory(folderID *string) (*Directory, error) {
	if folderID != nil && *folderID == "" {
		folderID = nil
	}

	folders, err := s.folderRepo.ListFolders(folderID)
	if err != nil {
		return nil, err
	}

	fileFolderID := ""
	if folderID != nil {
		fileFolderID = *folderID
	}

	files, err := s.fileRepo.ListFiles(fileFolderID)
	if err != nil {
		return nil, err
	}

	return &Directory{
		Folders: folders,
		Files:   files,
	}, nil
}

func (s *Service) findChildFolder(parentID *string, name string) (*model.Folder, error) {
	folders, err := s.folderRepo.ListFolders(parentID)
	if err != nil {
		return nil, err
	}

	for i := range folders {
		if folders[i].Name == name {
			return &folders[i], nil
		}
	}

	return nil, errFolderNotFound
}

func splitFolderPath(folderPath string) ([]string, error) {
	normalized := strings.TrimSpace(folderPath)
	if normalized == "" {
		return nil, nil
	}

	normalized = strings.ReplaceAll(normalized, "\\", "/")
	cleaned := path.Clean(normalized)
	if cleaned == "." || cleaned == "/" {
		return nil, nil
	}

	trimmed := strings.Trim(cleaned, "/")
	if trimmed == "" {
		return nil, nil
	}

	parts := strings.Split(trimmed, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return nil, fmt.Errorf("invalid folder path: %s", folderPath)
		}
	}

	return parts, nil
}

func rootFolder() *model.Folder {
	return &model.Folder{
		ID:   "",
		Name: "/",
	}
}
