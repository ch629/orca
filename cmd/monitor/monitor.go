package monitor

import (
	"context"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "monitor",
	Short: "Commands to monitor service status",
}

func init() {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the monitoring process",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Run this detached
			NewStart().Run(context.Background())
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops the monitoring process",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			NewStop().Run()
		},
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists current monitors",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			NewList().Run()
		},
	}

	var addFlags AddFlags
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new monitor",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			NewAdd().Run(addFlags)
		},
	}
	addFlags.SetupFlags(addCmd)

	removeCmd := &cobra.Command{
		Use:     "remove",
		Short:   "Removes a monitor",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"rm"},
		Run: func(cmd *cobra.Command, args []string) {
			NewRemove().Run(args[0])
		},
	}

	Cmd.AddCommand(startCmd, stopCmd, listCmd, addCmd, removeCmd)
}
