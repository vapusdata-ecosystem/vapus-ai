package cmd

import (
	"github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

const (
	getContextCmd     = "get-contexts"
	addContextCmd     = "add-context"
	currentContextCmd = "current-context"
	useContextCmd     = "use-context"
)

func NewConfigCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   pkg.ConfigResource,
		Short: "This command is parent command for all the config related operations.",
		Long: `This command is parent command for all the config related operations.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				plclient.MasterGlobals.Logger.Info().Msgf("No action provided for this command, use this command '" + pkg.ConfigResource + " --help' for the list of available commands")
			}
		},
	}
	cmd.AddCommand(NewGetContextsCmd(), NewAddContextCmd(), NewContextCurrentCmd(), NewUseContextCmd())
	return cmd
}
