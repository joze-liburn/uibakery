package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/joze-liburn/uibakery/shopify"
)

var (
	shopifyCmd = &cobra.Command{
		Use:   "shopify",
		Short: "Operations on Shopify",
		Long:  `Operations on Shopify.`,
	}

	shopifyListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists companies",
		Long:  `List companies ids; all or ones updated after a given date.`,
		Run:   shopifyListRun,
	}
)

func init() {
	shopifyCmd.AddCommand(shopifyListCmd)
	shopifyListCmd.Flags().Uint("limit", 10000, "Maximal number of records to fetch.")
	shopifyListCmd.Flags().Time("after", time.Now().Add(time.Hour), []string{"2006-01-02", time.RFC3339}, "Fetch only companies updated after this time.")

}

func shopifyListRun(cmd *cobra.Command, args []string) {
	host := viper.GetString("shp-hostname")
	if host == "" {
		fmt.Fprintf(os.Stderr, "Spotify address missing.")
		return
	}
	scrt := viper.GetString("shp-secret")
	client := shopify.New(host, scrt)

	limit, err := cmd.Flags().GetUint("limit")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing parameter limit: %s\n", err)
		return
	}
	var after *time.Time
	if cmd.Flags().Lookup("after") != nil {
		t, err := cmd.Flags().GetTime("after")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing parameter after: %s\n", err)
		}
		after = &t
	}

	stream := client.StreamCompaniesIds(10, limit, after)
	count := 0
	for nodeerr := range stream {
		if nodeerr.GetError() != nil {
			fmt.Printf("ERROR: %s\n", nodeerr.GetError())
			return
		}
		fmt.Printf("%3d: %s\n", count, nodeerr.(shopify.CompanyError).Company.Id)
		count++
	}
}
