package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "help",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.AddCommand(serveCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Running error: ", err)
		panic(err)
	}
}
