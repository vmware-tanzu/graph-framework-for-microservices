package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.eng.vmware.com/nexus/cli/pkg/common"
)

func Version(cmd *cobra.Command, args []string) error {
	fmt.Printf("CLI: %s\n", common.VERSION)
	fmt.Printf("BUILT: %s\n", common.BUILT)
	fmt.Printf("GIT_BRANCH: %s\n", common.GIT_BRANCH)
	fmt.Printf("GIT_COMMIT: %s\n", common.GIT_COMMIT)

	return nil
}
