package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/joze-liburn/uibakery/zendesk"
)

var (
	zendeskCmd = &cobra.Command{
		Use:   "zendesk",
		Short: "Operations on ZenDesk",
		Long:  `Operations on ZenDesk.`,
	}

	zendeskListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists companies",
		Long:  `List companies ids; all or ones updated after a given date.`,
		Run:   zendeskListRun,
	}
)

func init() {
	zendeskCmd.AddCommand(zendeskListCmd)
	zendeskListCmd.Flags().Uint("limit", 10000, "Maximal number of records to fetch.")
}

func zendeskListRun(cmd *cobra.Command, args []string) {
	host := viper.GetString("zen-hostname")
	if host == "" {
		fmt.Fprintf(os.Stderr, "Spotify address missing.")
		return
	}
	scrt := viper.GetString("zen-secret")
	client := zendesk.NewZendesk(host, scrt)

	limit, err := cmd.Flags().GetUint("limit")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing parameter limit: %s\n", err)
		return
	}

	stream := client.StreamOrganizations(10, limit)
	count := 0
	for organizationerr := range stream {
		if organizationerr.Err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", organizationerr.Err)
			return
		}
		fmt.Printf("%3d: %d %s\n", count, organizationerr.Organization.Id, organizationerr.Organization.ExternalId)
		count++
	}
}
