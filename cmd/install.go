package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install golang compiler",
	Long:  `安装 golang 开发环境`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

type installer struct {
	cmd *cobra.Command
}

func NewInstaller() Factory {
	var root = &installer{}
	root.init()
	return root
}

func (r *installer) init() {
	if r.cmd == nil {
		r.cmd = installCmd
	}
	if r.cmd.Run == nil {
		r.cmd.Run = r.run
	}
}

func (r *installer) run(cmd *cobra.Command, args []string) {
	log.Printf("args - %v \n", args)
}

func (r *installer) GetCmd() *cobra.Command {
	return r.cmd
}
