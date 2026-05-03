package vfs

// Service é um agregador dos services de folder e file
// Mantido para compatibilidade com código existente
type Service struct {
	*FolderService
	*FileService
}

func NewService(
	fileRepo fileRepository,
	folderRepo folderRepository,
) *Service {
	return &Service{
		FolderService: NewFolderService(folderRepo),
		FileService:   NewFileService(fileRepo),
	}
}

// ListDirectory combina folders e files de um diretório
func (s *Service) ListDirectory(folderID *string) (*Directory, error) {
	if folderID != nil && *folderID == "" {
		folderID = nil
	}

	folders, err := s.FolderService.ListFolders(folderID)
	if err != nil {
		return nil, err
	}

	fileFolderID := ""
	if folderID != nil {
		fileFolderID = *folderID
	}

	files, err := s.FileService.ListFiles(fileFolderID)
	if err != nil {
		return nil, err
	}

	return &Directory{
		Folders: folders,
		Files:   files,
	}, nil
}
