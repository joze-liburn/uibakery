package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

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
		Long: `Releases claimed record in the queue. Mandatory --claim-id flag
specifies the claim to release.`,
		Run: queueUnclaimRun,
	}

	queueCountCmd = &cobra.Command{
		Use:   "count",
		Short: "Counts records in the queue",
		Long: `Counts claimed and unclaimed records in the queue by status and
destination.`,
		Run: queueCountRun,
	}

	queueListClaimsCmd = &cobra.Command{
		Use:   "list-claims",
		Short: "Displays claims in the queue.",
		Long: `Displays claims in the queue, along with the number of records claimed.
Use --status to filter by status (exact match, case sensitive).`,
		Run: queueListClaimsRun,
	}
)

func init() {
	queueCmd.AddCommand(queueClaimCmd)
	queueClaimCmd.Flags().Uint("size", 10, "Up to this many records will be claimed.")
	viper.BindPFlag("size", queueClaimCmd.Flags().Lookup("size"))

	queueCmd.AddCommand(queueUnclaimCmd)
	queueUnclaimCmd.Flags().String("claim-id", "", "Id of the claim.")
	queueUnclaimCmd.MarkFlagRequired("claim-id")

	queueCmd.AddCommand(queueCountCmd)

	queueCmd.AddCommand(queueListClaimsCmd)
	queueListClaimsCmd.Flags().String("status", "", "Claims for this status only.")
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

func queueCountRun(cmd *cobra.Command, args []string) {
	db := &lbqueue.LbDb{}
	if err := db.Open(viper.GetString("username"), viper.GetString("secret"), viper.GetString("hostname"), viper.GetUint("port"), viper.GetString("database")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	stats, err := db.GetQueueCounts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetQueueStats() retured an error %s\n", err)
		return
	}
	printTable(stats)
}

func key(row string, col string) string {
	return row + "|" + col
}

func printTable(stats []lbqueue.QueueCount) {
	const (
		dataW = 8
	)
	type tabData struct {
		claimed   uint
		unclaimed uint
	}

	rows := []string{}
	cols := []string{}
	data := map[string]tabData{}

	for _, stat := range stats {
		rows = append(rows, stat.SubmissionStatus)
		cols = append(cols, stat.DestinationName)
		if stat.Claimed {
			data[key(stat.SubmissionStatus, stat.DestinationName)] = tabData{claimed: stat.Count}
		} else {
			data[key(stat.SubmissionStatus, stat.DestinationName)] = tabData{unclaimed: stat.Count}
		}
	}
	if len(data) == 0 {
		return
	}
	slices.Sort(rows)
	slices.Sort(cols)
	rows = slices.Compact(rows)
	cols = append(slices.Compact(cols), "total")
	col1w := 10
	for _, row := range rows {
		if col1w < len(row) {
			col1w = len(row)
		}
	}

	// Header row
	fmt.Printf("%*s | ", col1w, " ")
	for _, col := range cols {
		fmt.Printf("%*s | ", -(2*dataW + 3), col)
	}
	// Subheader row
	fmt.Printf("\n%*s | ", col1w, " ")
	for range cols {
		fmt.Printf("%*s : %*s | ", -dataW, "claimed", -dataW, "free")
	}
	// Separator
	fmt.Printf("\n%s\n", strings.Repeat("-", col1w+3+2*len(cols)*(dataW+3)-1))
	// Data rows
	cols = cols[:len(cols)-1]
	for _, row := range rows {
		var (
			claimT uint
			freeT  uint
		)
		fmt.Printf("%*s | ", -col1w, row)
		for _, col := range cols {
			ds := row + "|" + col
			c, ok := data[ds]
			if ok {
				fmt.Printf("%*d : %*d | ", dataW, c.claimed, dataW, c.unclaimed)
				claimT += c.claimed
				freeT += c.unclaimed
			} else {
				fmt.Printf("%*s : %*s | ", dataW, " ", dataW, " ")
			}
		}
		fmt.Printf("%*d : %*d |\n", dataW, claimT, dataW, freeT)
	}
}

func queueListClaimsRun(cmd *cobra.Command, args []string) {
	db := &lbqueue.LbDb{}
	if err := db.Open(viper.GetString("db-username"), viper.GetString("db-secret"), viper.GetString("db-hostname"), viper.GetUint("db-port"), viper.GetString("db-database")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	var status *string

	if cmd.Flags().Changed("status") {
		stat, _ := cmd.Flags().GetString("status")
		status = &stat
	}
	claims, err := db.ListClaims(status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ListClaims() retured an error %s\n", err)
		return
	}
	if len(claims) == 0 {
		fmt.Fprintf(os.Stderr, "No claims found\n")
		return
	}
	for _, claim := range claims {
		fmt.Printf("%s %6d\n", claim.Id, claim.Count)
	}
}
