package cmd

import (
	"fmt"
	"os"
	"strings"

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
	vpr := viper.New()
	initViper(vpr)

	for _, flag := range []struct {
		name   string
		defval any
		usage  string
		config string
	}{
		{name: "db-hostname", defval: "", usage: "hostname (localhost for SQL Auth Proxy)", config: "lbdatabase.hostname"},
		{name: "db-port", defval: uint16(5432), usage: "port number (5432 works for default SQL Auth Proxy)", config: "lbdatabase.port"},
		{name: "db-database", defval: "lightburn", usage: "Default database", config: "lbdatabase.database"},
		{name: "db-username", defval: "lb_ap_uibakery", usage: "User name (lb_ap_uibakery)", config: "lbdatabase.username"},
		{name: "db-secret", defval: "", usage: "Password", config: "lbdatabase.secret"},
		{name: "zen-hostname", defval: "lightburnsoftware.zendesk.com", usage: "ZenDFesk API endpoint", config: "zendesk.hostname"},
		{name: "zen-secret", defval: "", usage: "Access token for ZenDFesk API", config: "zendesk.secret"},
		{name: "shp-hostname", defval: "lightburn-software-llc.myshopify.com", usage: "Shopify API endpoint", config: "shopify.hostname"},
		{name: "shp-secret", defval: "", usage: "Shopify API secret", config: "shopify.secret"},
		{name: "crp-hostname", defval: "api.cryptlex.com", usage: "Cryptlex API hostname", config: "cryptlex.hostname"},
		{name: "crp-secret", defval: "", usage: "Cryptlex secret", config: "cryptlex.secret"},
		{name: "crp-version", defval: uint(3), usage: "Cryptlex API version", config: "cryptlex.version"},
	} {
		switch tval := flag.defval.(type) {
		case string:
			rootCmd.PersistentFlags().String(flag.name, tval, flag.usage)
		case uint16:
			rootCmd.PersistentFlags().Uint16(flag.name, tval, flag.usage)
		case uint:
			rootCmd.PersistentFlags().Uint(flag.name, tval, flag.usage)
		default:
			fmt.Fprintf(os.Stderr, "%s: type %T not implemented", flag.name, tval)
		}
		flg := rootCmd.PersistentFlags().Lookup(flag.name)
		if flg == nil {
			continue
		}
		if err := vpr.BindPFlag(flag.name, flg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}
		if flag.config == "" {
			continue
		}
		// This part passes a value as set by env or config file (that is, from
		// viper) to command line flag (cobra) _if_ command line value is not
		// the default one.
		if !flg.Changed && vpr.IsSet(flag.config) {
			val := vpr.Get(flag.config)
			if err := rootCmd.PersistentFlags().Set(flag.name, fmt.Sprintf("%v", val)); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}
			vpr.Set(flag.name, fmt.Sprintf("%v", val))
		}
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(shopifyCmd)
	rootCmd.AddCommand(zendeskCmd)
	rootCmd.AddCommand(cryptlexCmd)
}

func initConfig() {

}

func initViper(v *viper.Viper) {
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".uibakery" (without extension).
		v.AddConfigPath(home)
		v.AddConfigPath(".")
		v.SetConfigName(".uibakery")
		v.SetConfigType("json")
	}
	v.ReadInConfig()

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of uibakery",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uibakery data sync server v0.1")
	},
}
