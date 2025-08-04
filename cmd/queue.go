package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/joze-liburn/uibakery/lbqueue"
)

var (
	queueCmd = &cobra.Command{
		Use:   "queue",
		Short: "Operations over the queue",
		Long:  `Operations over the queue.`,
	}

	queueClaimCmd = &cobra.Command{
		Use:   "claim",
		Short: "Claims record from the queue",
		Long:  `Claims record from the queue.`,
		Run:   queueClaimRun,
	}

	queueUnclaimCmd = &cobra.Command{
		Use:   "unclaim",
		Short: "Releases claimed record in the queue",
		Long:  `Releases claimed record in the queue.`,
		Run:   queueUnclaimRun,
	}
)

func init() {
	queueCmd.AddCommand(queueClaimCmd)
	queueClaimCmd.Flags().Uint("size", 10, "Up to this many records will be claimed.")
	viper.BindPFlag("size", queueClaimCmd.Flags().Lookup("size"))

	queueCmd.AddCommand(queueUnclaimCmd)
	queueUnclaimCmd.Flags().String("claim-id", "", "Id of the claim.")
	queueUnclaimCmd.MarkFlagRequired("claim-id")
	viper.BindPFlag("claim-id", queueClaimCmd.Flags().Lookup("claim-id"))
}

func queueClaimRun(cmd *cobra.Command, args []string) {
	db := &lbqueue.LbDb{}
	if err := db.Open(viper.GetString("username"), viper.GetString("secret"), viper.GetString("hostname"), viper.GetUint("port"), viper.GetString("database")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	size, err := cmd.Flags().GetUint("size")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ClaimRecords() retured an error %s\n", err)
		return
	}
	g, c, err := db.ClaimRecords(size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ClaimRecords() retured an error %s\n", err)
		return
	}
	fmt.Printf("Claim %q of %d records\n", g, c)
}

func queueUnclaimRun(cmd *cobra.Command, args []string) {
	db := &lbqueue.LbDb{}
	if err := db.Open(viper.GetString("username"), viper.GetString("secret"), viper.GetString("hostname"), viper.GetUint("port"), viper.GetString("database")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	claim, err := cmd.Flags().GetString("claim-id")
	if err != nil {
		fmt.Fprintf(os.Stderr, "UnclaimRecords(%q) retured an error %s\n", claim, err)
		return
	}
	c, err := db.UnclaimRecords(claim)
	if err != nil {
		fmt.Fprintf(os.Stderr, "UnclaimRecords(%q) retured an error %s\n", claim, err)
		return
	}
	fmt.Printf("%d records released from claim %q\n", c, claim)
}
