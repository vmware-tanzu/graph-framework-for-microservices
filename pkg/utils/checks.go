package utils

import "github.com/spf13/cobra"

const (
	EnableDebugFlag     = "debug"
	ListPrereqFlag      = "list-prereq"
	SkipPrereqCheckFlag = "skip-prereq-check"
)

func IsDebug(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(EnableDebugFlag).Changed
}

func VerifyAll(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup("all").Changed
}

func ListPrereq(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(ListPrereqFlag).Changed
}

func SkipPrereqCheck(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup(SkipPrereqCheckFlag).Changed
}
