package repository

import (
	"database/sql"
	"drivepro/internal/storage/model"
	"time"

	"github.com/google/uuid"
)

type FolderRepository struct {
	DB *sql.DB
}

func NewFolderRepository(db *sql.DB) *FolderRepository {
	return &FolderRepository{DB: db}
}

func (r *FolderRepository) CreateFolder(name string, parentID *string) (*model.Folder, error) {
	folder := &model.Folder{
		ID:        uuid.NewString(),
		Name:      name,
		ParentID:  parentID,
		CreatedAt: time.Now().Unix(),
	}

	_, err := r.DB.Exec(`
		INSERT INTO folders (id, name, parent_id, created_at)
		VALUES (?, ?, ?, ?)
	`, folder.ID, folder.Name, folder.ParentID, folder.CreatedAt)

	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *FolderRepository) GetFolderByID(id string) (*model.Folder, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, parent_id, created_at
		FROM folders
		WHERE id = ?
	`, id)

	var folder model.Folder
	var parentID sql.NullString

	err := row.Scan(&folder.ID, &folder.Name, &parentID, &folder.CreatedAt)
	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		folder.ParentID = &parentID.String
	}

	return &folder, nil
}

func (r *FolderRepository) ListFolders(parentID *string) ([]model.Folder, error) {
	var rows *sql.Rows
	var err error

	if parentID == nil {
		rows, err = r.DB.Query(`
			SELECT id, name, parent_id, created_at
			FROM folders
			WHERE parent_id IS NULL
		`)
	} else {
		rows, err = r.DB.Query(`
			SELECT id, name, parent_id, created_at
			FROM folders
			WHERE parent_id = ?
		`, *parentID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []model.Folder

	for rows.Next() {
		var folder model.Folder
		var parent sql.NullString

		err := rows.Scan(&folder.ID, &folder.Name, &parent, &folder.CreatedAt)
		if err != nil {
			return nil, err
		}

		if parent.Valid {
			folder.ParentID = &parent.String
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *FolderRepository) DeleteFolder(id string) error {
	_, err := r.DB.Exec(`DELETE FROM folders WHERE id = ?`, id)
	return err
}
