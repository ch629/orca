package main

import (
	"github.com/ch629/orca/cmd/monitor"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "orca",
}

func init() {
	rootCmd.AddCommand(monitor.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
