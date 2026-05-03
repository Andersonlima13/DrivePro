package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "List directory contents",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		folder, err := service.GetFolderByPath(path)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		dir, err := service.ListDirectory(&folder.ID)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, f := range dir.Folders {
			fmt.Println("[DIR]", f.Name)
		}

		for _, f := range dir.Files {
			fmt.Println("[FILE]", f.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
