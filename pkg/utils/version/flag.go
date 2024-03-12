package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GitCommit string
	BuildDate string
	GoVersion string
	Platform  string
	Version   string
)

func AddFlags(cmd *cobra.Command) {
	v := fmt.Sprintf("Version: \"%s\"\n", Version) +
		fmt.Sprintf("GitCommit: \"%s\"\n", GitCommit) +
		fmt.Sprintf("BuildDate: \"%s\"\n", BuildDate) +
		fmt.Sprintf("GoVersion: \"%s\"\n", GoVersion) +
		fmt.Sprintf("Platform: \"%s\"\n", Platform)
	cmd.Version = v
	cmd.SetVersionTemplate(`{{printf "%s" .Version}}`)
}
