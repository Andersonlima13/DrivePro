package storage

type FileStatus string

const (
	FileStatusPending FileStatus = "pendente"
	FileStatusSyncing FileStatus = "sincronizando"
	FileStatusSynced  FileStatus = "sincronizado"
)
