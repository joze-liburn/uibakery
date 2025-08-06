package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "uibakery",
		Short: "A data synchromnization server.",
		Long:  `uibakery is used to synchronize data from Shopify to ZenDesk and Criplex`,
	}
)

// Execute executes the root command.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.uibakery.json)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	viper.SetDefault("db-hostname", "localhost")
	viper.SetDefault("db-port", 5431)
	viper.SetDefault("db-username", "lb_ap_uibakery")
	viper.SetDefault("shp-hostname", "lightburn-software-llc.myshopify.com")

	for _, flag := range []struct {
		name   string
		defval any
		usage  string
	}{
		{name: "db-hostname", defval: "", usage: "hostname (localhost for SQL Auth Proxy)"},
		{name: "db-port", defval: uint16(5432), usage: "port number (5432 works for default SQL Auth Proxy)"},
		{name: "db-database", defval: "lightburn", usage: "Default database"},
		{name: "db-username", defval: "lb_ap_uibakery", usage: "User name (lb_ap_uibakery)"},
		{name: "db-secret", defval: "", usage: "Password"},
		{name: "zen-hostname", defval: "https://lightburnsoftware.zendesk.com/api/v2/", usage: "ZenDFesk API endpoint"},
		{name: "zen-secret", defval: "", usage: "Access token for ZenDFesk API"},
		{name: "shp-hostname", defval: "lightburn-software-llc.myshopify.com", usage: "Shopify API endpoint"},
		{name: "shp-secret", defval: "", usage: "Shopify API secret"},
	} {
		switch tval := flag.defval.(type) {
		case string:
			rootCmd.PersistentFlags().String(flag.name, tval, flag.usage)
		case uint16:
			rootCmd.PersistentFlags().Uint16(flag.name, tval, flag.usage)
		default:
			fmt.Fprintf(os.Stderr, "%s: type %T not implemented", flag.name, tval)
		}
		if err := viper.BindPFlag(flag.name, rootCmd.PersistentFlags().Lookup(flag.name)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(shopifyCmd)
	rootCmd.AddCommand(zendeskCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".uibakery")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of uibakery",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uibakery data sync server v0.1")
	},
}
