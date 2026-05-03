package vfs

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"drivepro/internal/storage/model"
)

var errFolderNotFound = errors.New("folder not found")

type folderRepository interface {
	CreateFolder(name string, parentID *string) (*model.Folder, error)
	ListFolders(parentID *string) ([]model.Folder, error)
}

type FolderService struct {
	folderRepo folderRepository
}

func NewFolderService(folderRepo folderRepository) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
	}
}

func (s *FolderService) CreateFolderByPath(folderPath string) error {
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

func (s *FolderService) GetFolderByPath(folderPath string) (*model.Folder, error) {
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

func (s *FolderService) ListFolders(parentID *string) ([]model.Folder, error) {
	return s.folderRepo.ListFolders(parentID)
}

func (s *FolderService) findChildFolder(parentID *string, name string) (*model.Folder, error) {
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
