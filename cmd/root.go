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

	RootCmd = &cobra.Command{
		Use:   "uibakery",
		Short: "A data synchromnization server.",
		Long:  `uibakery is used to synchronize data from Shopify to ZenDesk and Criplex`,
	}
)

// Execute executes the root command.
func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.uibakery.json)")
	RootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))

	viper.SetDefault("hostname", "localhost")
	viper.SetDefault("port", 5431)
	viper.SetDefault("username", "lb_ap_uibakery")

	queueCmd.PersistentFlags().String("db-hostname", "", "hostname (localhost for SQL Auth Proxy)")
	queueCmd.PersistentFlags().Uint16("db-port", 5432, "port number (5432 works for default SQL Auth Proxy)")
	queueCmd.PersistentFlags().String("db-database", "lightburn", "Default database")
	queueCmd.PersistentFlags().String("db-username", "lb_ap_uibakery", "User name (lb_ap_uibakery)")
	queueCmd.PersistentFlags().String("db-secret", "", "Password")
	queueCmd.PersistentFlags().String("zen-hostname", "https://lightburnsoftware.zendesk.com/api/v2/", "ZenDFesk API endpoint")
	queueCmd.PersistentFlags().String("zen-secret", "", "Access token for ZenDFesk API")

	for _, k := range []string{"db-hostname", "db-port", "db-database", "db-username", "db-secret", "zen-hostname", "zen-secret"} {
		if err := viper.BindPFlag(k, queueCmd.PersistentFlags().Lookup(k)); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(serverCmd)
	RootCmd.AddCommand(queueCmd)
	RootCmd.AddCommand(shopifyCmd)
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
