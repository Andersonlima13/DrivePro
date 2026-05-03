package main

import (
	"drivepro/internal/storage"
	"drivepro/internal/storage/repository"
	"drivepro/internal/vfs"
)

var service *vfs.Service

func initService() {
	storage.InitDB("drivepro.db")

	fileRepo := repository.NewFileRepository(storage.DB)
	folderRepo := repository.NewFolderRepository(storage.DB)

	service = vfs.NewService(fileRepo, folderRepo)
}
