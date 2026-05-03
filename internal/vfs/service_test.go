package vfs

import (
	"fmt"
	"testing"

	"drivepro/internal/storage/model"
)

func TestServiceFolderPathAndDirectoryListing(t *testing.T) {
	service, fileRepo := newTestService()

	if err := service.CreateFolderByPath("/docs/projects"); err != nil {
		t.Fatalf("CreateFolderByPath() error = %v", err)
	}

	root, err := service.GetFolderByPath("/")
	if err != nil {
		t.Fatalf("GetFolderByPath(\"/\") error = %v", err)
	}
	if root.ID != "" || root.Name != "/" {
		t.Fatalf("root folder = %#v, want virtual root", root)
	}

	docs, err := service.GetFolderByPath("/docs")
	if err != nil {
		t.Fatalf("GetFolderByPath(\"/docs\") error = %v", err)
	}

	projects, err := service.GetFolderByPath("\\docs\\projects")
	if err != nil {
		t.Fatalf("GetFolderByPath(\"\\\\docs\\\\projects\") error = %v", err)
	}
	if projects.ParentID == nil || *projects.ParentID != docs.ID {
		t.Fatalf("projects parent = %v, want %s", projects.ParentID, docs.ID)
	}

	fileRepo.files = append(fileRepo.files, model.File{
		ID:       "file-1",
		Name:     "notes.txt",
		Size:     12,
		Type:     "text/plain",
		FolderID: docs.ID,
	})

	dir, err := service.ListDirectory(&docs.ID)
	if err != nil {
		t.Fatalf("ListDirectory() error = %v", err)
	}
	if len(dir.Folders) != 1 || dir.Folders[0].Name != "projects" {
		t.Fatalf("folders = %#v, want projects folder", dir.Folders)
	}
	if len(dir.Files) != 1 || dir.Files[0].Name != "notes.txt" {
		t.Fatalf("files = %#v, want notes.txt file", dir.Files)
	}

	rootDir, err := service.ListDirectory(nil)
	if err != nil {
		t.Fatalf("ListDirectory(nil) error = %v", err)
	}
	if len(rootDir.Folders) != 1 || rootDir.Folders[0].Name != "docs" {
		t.Fatalf("root folders = %#v, want docs folder", rootDir.Folders)
	}
}

func newTestService() (*Service, *fakeFileRepository) {
	fileRepo := &fakeFileRepository{}
	folderRepo := &fakeFolderRepository{}

	return NewService(fileRepo, folderRepo), fileRepo
}

type fakeFileRepository struct {
	files []model.File
}

func (r *fakeFileRepository) ListFiles(folderID string) ([]model.File, error) {
	var files []model.File
	for _, file := range r.files {
		if file.FolderID == folderID {
			files = append(files, file)
		}
	}

	return files, nil
}

type fakeFolderRepository struct {
	nextID  int
	folders []model.Folder
}

func (r *fakeFolderRepository) CreateFolder(name string, parentID *string) (*model.Folder, error) {
	r.nextID++

	folder := model.Folder{
		ID:   fmt.Sprintf("folder-%d", r.nextID),
		Name: name,
	}
	if parentID != nil {
		parent := *parentID
		folder.ParentID = &parent
	}

	r.folders = append(r.folders, folder)

	return &folder, nil
}

func (r *fakeFolderRepository) ListFolders(parentID *string) ([]model.Folder, error) {
	var folders []model.Folder
	for _, folder := range r.folders {
		if sameParent(folder.ParentID, parentID) {
			folders = append(folders, folder)
		}
	}

	return folders, nil
}

func sameParent(a *string, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	return *a == *b
}
