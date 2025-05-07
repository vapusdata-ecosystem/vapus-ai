package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

func globalPostRun(cmd *cobra.Command, args []string) {
	fmt.Println("Post-run hook executed for command:", cmd.Name())
	// Setting a sample value in GlobalVar to demonstrate
	GlobalVar = fmt.Sprintf("Value set by %s", cmd.Name())
}

func NewConnectCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.ConnectOps,
		Short: "This command is to connect with current context.",
		Long:  `This command is to connect with current context.`,
		Run: func(cmd *cobra.Command, args []string) {
			plclient.MasterGlobals.Logger.Info().Msg("Connecting to current context")
			currentContext := viper.GetString(currentContextKey)
			if currentContext == "" {
				cobra.CheckErr(pkg.ErrNoCurrentContext)
			}
			GlobalVar = "Data from cc1"
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			setGlobalPersist()
			globalPostRun(cmd, args)
			initCurrentContextInstance()
			return nil
		},
	}
	return cmd
}
