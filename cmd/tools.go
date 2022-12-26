package cmd

import (
	"github.com/spf13/cobra"
	tools2 "github.com/weblfe/dragon/pkg/tools"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	if r.cmd.PreRun == nil {
		r.cmd.PreRun = r.Check
	}
}

func (r *tools) Check(_ *cobra.Command, _ []string) {
	root := os.Getenv("GOROOT")
	_path := os.Getenv("GOPATH")
	if !strings.Contains(root, string(filepath.Separator)) {
		root = strings.ReplaceAll(root, "/", string(filepath.Separator))
	}
	if !strings.Contains(_path, string(filepath.Separator)) {
		_path = strings.ReplaceAll(_path, "/", string(filepath.Separator))
	}
	log.Printf("GOROOT: %s\n", root)
	log.Printf("GOPATH: %s\n", _path)
	if root == _path && root != "" {
		_path = filepath.Join(filepath.Dir(root), "go_path")
		if _, err := os.Stat(_path); os.IsNotExist(err) {
			err = os.MkdirAll(_path, os.ModePerm)
			panic(err)
		}
		_ = os.Setenv("GOPATH", _path)
	}
}

func (r *tools) parseArgs(args []string) map[string]string {
	var kvs = Kvs{}
	for _, v := range args {
		if !strings.Contains(v, "=") {
			kvs[v] = ""
			continue
		}
		values := strings.SplitN(v, "=", 2)
		kvs[values[0]] = values[1]
	}
	return kvs
}

func (r *tools) run(_ *cobra.Command, args []string) {
	log.Printf("args - %v \n", args)
	kvs := r.parseArgs(args)
	kvs["GOROOT"] = goRoot
	kvs["GOPATH"] = goPath
	devOps := tools2.NewDevops(kvs)
	// 罗列支持安装的工具列表
	if _, ok := kvs["lists"]; ok {
		devOps.Lists()
		return
	}
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
