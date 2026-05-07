package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		runSQLiteFolderSmokeTest()
		return
	}

	Execute()
}

func runSQLiteFolderSmokeTest() {
	initService()

	folders := []string{
		"/teste-sqlite-1",
		"/teste-sqlite-2",
	}

	for _, folderPath := range folders {
		ensureFolder(folderPath)
	}

	dir, err := service.ListDirectory(nil)
	if err != nil {
		exitWithError(err)
	}

	fmt.Println("Root folders in SQLite:")
	for _, folder := range dir.Folders {
		fmt.Println("[DIR]", folder.Name)
	}
}

func ensureFolder(folderPath string) {
	err := service.CreateFolderByPath(folderPath)
	if err == nil {
		fmt.Println("Folder created:", folderPath)
		return
	}

	if _, err := service.GetFolderByPath(folderPath); err == nil {
		fmt.Println("Folder already exists:", folderPath)
		return
	}

	exitWithError(fmt.Errorf("create folder %q: %w", folderPath, err))
}

func exitWithError(err error) {
	fmt.Println("Error:", err)
	os.Exit(1)
}
