package repository

import (
	"database/sql"
	"drivepro/internal/storage/model"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FileRepository struct {
	DB *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{DB: db}
}

func (r *FileRepository) CreateFile(name string, size int64, fileType string, folderID string) (*model.File, error) {
	file := &model.File{
		ID:        uuid.NewString(),
		Name:      name,
		Size:      size,
		Type:      fileType,
		FolderID:  folderID,
		Status:    model.FileStatusPending,
		CreatedAt: time.Now().Unix(),
	}

	_, err := r.DB.Exec(`
		INSERT INTO files (id, name, size, type, folder_id, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, file.ID, file.Name, file.Size, file.Type, file.FolderID, file.Status, file.CreatedAt)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r *FileRepository) GetFileByID(id string) (*model.File, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, size, type, folder_id, status, created_at
		FROM files
		WHERE id = ?
	`, id)

	var file model.File

	err := row.Scan(
		&file.ID,
		&file.Name,
		&file.Size,
		&file.Type,
		&file.FolderID,
		&file.Status,
		&file.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) ListFiles(folderID string) ([]model.File, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, size, type, folder_id, status, created_at
		FROM files
		WHERE folder_id = ?
	`, folderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.File

	for rows.Next() {
		var file model.File

		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Size,
			&file.Type,
			&file.FolderID,
			&file.Status,
			&file.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (r *FileRepository) UpdateFileStatus(id string, status model.FileStatus) error {
	_, err := r.DB.Exec(`
		UPDATE files SET status = ?
		WHERE id = ?
	`, status, id)

	return err
}

func (r *FileRepository) DeleteFile(id string) error {
	_, err := r.DB.Exec(`DELETE FROM files WHERE id = ?`, id)
	return err
}

func (r *FileRepository) MoveFile(fileID string, newFolderID string) error {
	res, err := r.DB.Exec(`
		UPDATE files SET folder_id = ?
		WHERE id = ?
	`, newFolderID, fileID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("file not found")
	}

	return nil
}
