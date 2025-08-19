package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/joze-liburn/uibakery/cryptlex"
)

var (
	cryptlexCmd = &cobra.Command{
		Use:   "cryptlex",
		Short: "Operations on Cryptlex",
		Long:  `Operations on Cryptlex.`,
	}

	cryptlexListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists resellers",
		Long: fmt.Sprintf(`List resellers; specify criteria to limit search.
Flags support either a value (--sort <column>) or condition (--id "eq 1").
Condition is composed of operator and operand, separated by a space, Use double
quotes (") or escape space with a backslash to pass space into argument.
Time formats: %s.
Again, eventual space must be escaped or the whole argument (together with the
operator) quoted.`, strings.Join(cryptlex.TimeFormats, ", ")),
		Run: cryptlexListRun,
	}
)

func init() {
	cryptlexCmd.AddCommand(cryptlexListCmd)
	addFlags(cryptlexCmd)
}

func addFlags(cmd *cobra.Command) {
	deft := time.Now().Add(-time.Hour)
	defs := deft.Format("ge 2006-01-02 15:04")
	cmd.Flags().Int("page", 1, "Page number to retrieve.")
	cmd.Flags().Int("limit", 100, "page length (1..100).")
	cmd.Flags().String("sort", "", "Order results on this attribute.")
	cmd.Flags().String("name", "", fmt.Sprintf(`(eg. --name "in Jane,Janet") Name of the reseller.
Operators: %s.`, strings.Join(cryptlex.OperatorsName, ",")))
	cmd.Flags().String("email", "", fmt.Sprintf(`(eg. --email eq\ jane@example.com) Reseller's email.
Operators: %s.`, strings.Join(cryptlex.OperatorsEmail, ",")))
	cmd.Flags().String("search", "", "Search string.")
	cmd.Flags().String("id", "", fmt.Sprintf(`(eg --id "ne 123") Reseller's id.
Operators: %s.`, strings.Join(cryptlex.OperatorsId, ",")))
	cmd.Flags().String("createdAt", "", fmt.Sprintf(`(eg. --createdAt "%s") searches by the creation time.
Operators: %s.`, defs, strings.Join(cryptlex.OperatorsCreatedAt, ",")))
	cmd.Flags().String("updatedAt", "", fmt.Sprintf(`(eg. --updatedAt "%s") searches by the last update time.
Operators: %s.`, defs, strings.Join(cryptlex.OperatorsUpdatedAt, ",")))
}

func client(cmd *cobra.Command) (*cryptlex.Cryptlex, error) {
	host := viper.GetString("crp-hostname")
	if host == "" {
		return nil, fmt.Errorf("Cryptlex address missing")
	}
	version, _ := cmd.Flags().GetUint("crp-version")
	scrt := viper.GetString("crp-secret")
	return cryptlex.NewCryptlex(version, host, scrt), nil
}

func parseFlags(cmd *cobra.Command) ([]cryptlex.CallOptions, error) {
	opts := []cryptlex.CallOptions{}
	for _, flg := range []struct {
		name string
		fnc  func(int) cryptlex.CallOptions
	}{
		{"page", cryptlex.WithPage},
		{"limit", cryptlex.WithPageSize},
	} {
		if f := cmd.Flags().Lookup(flg.name); f != nil && f.Changed {
			v, e := cmd.Flags().GetInt(flg.name)
			if e != nil {
				return nil, e
			}
			opts = append(opts, flg.fnc(v))
		}
	}
	if f := cmd.Flags().Lookup("sort"); f != nil && f.Changed {
		opts = append(opts, cryptlex.WithSort(f.Value.String()))
	}
	if f := cmd.Flags().Lookup("search"); f != nil && f.Changed {
		opts = append(opts, cryptlex.WithSearch(f.Value.String()))
	}
	for _, flg := range []struct {
		name string
		fnc  func(string) (cryptlex.CallOptions, error)
	}{
		{"id", cryptlex.WithIdFlag},
		{"name", cryptlex.WithNameFlag},
		{"email", cryptlex.WithEmailFlag},
		{"createdAt", cryptlex.WithCreatedAtFlag},
		{"updatedAt", cryptlex.WithUpdatedAtFlag},
	} {
		if f := cmd.Flags().Lookup(flg.name); f != nil && f.Changed {
			if f, e := flg.fnc(f.Value.String()); e != nil {
				return nil, e
			} else {
				opts = append(opts, f)
			}
		}
	}
	return opts, nil
}

func cryptlexListRun(cmd *cobra.Command, args []string) {
	client, err := client(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	opts, err := parseFlags(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	rsls, err := client.ListResellers(opts...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	count := 0
	for _, rsl := range rsls {
		fmt.Printf("%3d: %s %s %s\n", count, rsl.Id, rsl.Name, rsl.Email)
		count++
	}
}
