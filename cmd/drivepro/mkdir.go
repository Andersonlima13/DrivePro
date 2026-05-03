package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mkdirCmd = &cobra.Command{
	Use:   "mkdir [path]",
	Short: "Create a folder",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		err := service.CreateFolderByPath(path)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Folder created:", path)
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}
