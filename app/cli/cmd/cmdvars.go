package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	plclient "github.com/vapusdata-ecosystem/vapusdata/cli/internals/platform"
	pkg "github.com/vapusdata-ecosystem/vapusdata/cli/pkgs"
)

type PlatformInstanceClients struct {
	*plclient.VapusCtlClient
}

var (

	// RootCmd is the root command for vapusdt
	rootCmd                                                       *cobra.Command
	Logger                                                        zerolog.Logger
	currentIdToken, currentAccessToken, currentProductAccessToken string = "currentIdToken", "currentAccessToken", "currentProductAccessToken"
	GlobalVar                                                     string
)

var ignoreConnMap = map[string]bool{
	pkg.ConfigResource: true,
	pkg.ClearOps:       true,
	pkg.ExplainOps:     true,
	pkg.OperatorOps:    true,
	pkg.InstallerOps:   true,
}
