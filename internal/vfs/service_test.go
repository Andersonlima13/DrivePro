package vfs

import (
	"fmt"
	"testing"

	"drivepro/internal/storage/model"
)

func TestFolderServicePathAndDirectoryListing(t *testing.T) {
	folderService, fileService, fileRepo := newTestServices()

	if err := folderService.CreateFolderByPath("/docs/projects"); err != nil {
		t.Fatalf("CreateFolderByPath() error = %v", err)
	}

	root, err := folderService.GetFolderByPath("/")
	if err != nil {
		t.Fatalf("GetFolderByPath(\"/\") error = %v", err)
	}
	if root.ID != "" || root.Name != "/" {
		t.Fatalf("root folder = %#v, want virtual root", root)
	}

	docs, err := folderService.GetFolderByPath("/docs")
	if err != nil {
		t.Fatalf("GetFolderByPath(\"/docs\") error = %v", err)
	}

	projects, err := folderService.GetFolderByPath("\\docs\\projects")
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

	// Test listing folders in docs directory
	folders, err := folderService.ListFolders(&docs.ID)
	if err != nil {
		t.Fatalf("ListFolders() error = %v", err)
	}
	if len(folders) != 1 || folders[0].Name != "projects" {
		t.Fatalf("folders = %#v, want projects folder", folders)
	}

	// Test listing files in docs directory
	files, err := fileService.ListFiles(docs.ID)
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}
	if len(files) != 1 || files[0].Name != "notes.txt" {
		t.Fatalf("files = %#v, want notes.txt file", files)
	}

	// Test listing folders in root directory
	rootFolders, err := folderService.ListFolders(nil)
	if err != nil {
		t.Fatalf("ListFolders(nil) error = %v", err)
	}
	if len(rootFolders) != 1 || rootFolders[0].Name != "docs" {
		t.Fatalf("root folders = %#v, want docs folder", rootFolders)
	}
}

func newTestServices() (*FolderService, *FileService, *fakeFileRepository) {
	fileRepo := &fakeFileRepository{}
	folderRepo := &fakeFolderRepository{}

	return NewFolderService(folderRepo), NewFileService(fileRepo), fileRepo
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
