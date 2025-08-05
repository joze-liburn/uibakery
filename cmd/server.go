package cmd

import (
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "uibakery",
		Short: "Start synchronization server.",
		Long:  `Start synchronization server. Comand "never" returns.`,
	}
)
