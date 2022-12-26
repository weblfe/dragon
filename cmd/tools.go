package cmd

import (
	"github.com/spf13/cobra"
	tools2 "github.com/weblfe/dragon/pkg/tools"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "install golang tools",
	Long:  `安装 golang 相关开发工具`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
}

var (
	goRoot string
	goPath string
)

func init() {
	toolsCmd.Flags().StringVar(&goRoot, "go_root", "", "set GOROOT env")
	toolsCmd.Flags().StringVar(&goPath, "go_path", "", "set GOPATH env")
}

type tools struct {
	cmd *cobra.Command
}

type Kvs map[string]string

func NewTools() Factory {
	var root = &tools{}
	root.init()
	return root
}

func (r *tools) init() {
	if r.cmd == nil {
		r.cmd = toolsCmd
	}
	if r.cmd.Run == nil {
		r.cmd.Run = r.run
	}
}

func (r *tools) run(_ *cobra.Command, args []string) {
	log.Printf("args - %v \n", args)
	devOps := tools2.NewDevops(Kvs{"GOROOT": goRoot, "GOPATH": goPath})
	stdout, stderr, err := devOps.Exec()
	if err != nil {
		log.Printf("[ERROR] %v\n", err)
	}
	if stderr.Len() > 0 {
		log.Printf("[ERROR] %v\n", stderr.String())
	}
	if stdout.Len() > 0 {
		log.Printf("%v\n", stdout.String())
	}
}

func (r *tools) GetCmd() *cobra.Command {
	return r.cmd
}
