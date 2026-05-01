package model

type FileStatus string

const (
	FileStatusPending FileStatus = "PENDING"
	FileStatusSyncing FileStatus = "SYNCING"
	FileStatusSynced  FileStatus = "SYNCED"
	FileStatusError   FileStatus = "ERROR"
)
