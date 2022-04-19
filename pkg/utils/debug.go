package utils

import "github.com/spf13/cobra"

func IsDebug(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup("debug").Changed
}
