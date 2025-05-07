package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	pkg "github.com/vapusdata-ecosystem/vapusai/cli/pkgs"
)

// authCmd represents the auth command
func NewAuthCmd() *cobra.Command {

	authCmd := &cobra.Command{
		Use:   pkg.AuthOps,
		Short: "Login to the current data marketplace instance.",
		Long:  `This command is used to login to the VapusData platform`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("auth called")
		},
	}
	authCmd.AddCommand(NewAuthLoginCmd(), NewOrganizationAuthCmd(), NewDataProductAuthCmd())
	return authCmd
}

// func init() {
// 	rootCmd.AddCommand(authCmd)

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
