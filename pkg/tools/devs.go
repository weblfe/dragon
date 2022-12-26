package tools

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Devops struct {
	args  map[string]string
	tools map[Name]PkgSrc
}

type Name string
type PkgSrc string

func (n Name) String() string {
	return string(n)
}

func (pkg PkgSrc) String() string {
	return string(pkg)
}

func (pkg PkgSrc) Join(ver string) PkgSrc {
	return PkgSrc(pkg.String() + "@" + ver)
}

func NewDevops(args map[string]string) *Devops {
	var dev = &Devops{args: args}
	dev.init()
	return dev
}

func (dev *Devops) init() {
	dev.tools = createDefaultTools()
	dev.loadEnv()
}

func (dev *Devops) loadEnv() {
	root := os.Getenv("GOROOT")
	goPath := os.Getenv("GOPATH")
	log.Printf("GOROOT:%s \n", root)
	log.Printf("GOPATH:%s \n", goPath)
}

func (dev *Devops) Exec() (stdOut, stdErr bytes.Buffer, err error) {
	var (
		cmd            *exec.Cmd
		stdout, stderr bytes.Buffer
		args           = dev.GetArgs()
		name           = dev.GetCommand()
	)
	for _, v := range args {
		cmd = exec.Command(name, v)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err = cmd.Run(); err != nil {
			return stdout, stderr, err
		}
		log.Printf("%v", string(stdout.Bytes()))
		stdout.Reset()
		stdErr.Reset()
	}
	return stdout, stderr, nil
}

func (dev *Devops) GetCommand() string {
	root, ok := dev.args["GOROOT"]
	if !ok || root == "" {
		return "go get -u -v"
	}
	return dev.getGoBin(root)
}

func (dev *Devops) getGoBin(root string) string {
	var bin = filepath.Join(root, "bin", "go")
	if runtime.GOOS == "windows" {
		bin = filepath.Join(root, "bin", "go.exe")
	}
	return fmt.Sprintf("%s get -u -v", bin)
}

func (dev *Devops) GetArgs() []string {
	var (
		args  []string
		lists = dev.tools
	)
	tools, ok := dev.args["tools"]
	if !ok {
		for _, v := range lists {
			args = append(args, v.String())
		}
		return args
	}
	for _, v := range strings.Split(tools, ",") {
		var name, ver string
		if !strings.Contains(v, "@") {
			name = v
		} else {
			values := strings.Split(v, "@")
			name = values[0]
			if len(values) > 1 {
				ver = values[1]
			}
		}
		pkg, ok1 := lists[Name(name)]
		if !ok1 {
			continue
		}
		if ver == "" {
			args = append(args, pkg.String())
		} else {
			args = append(args, pkg.Join(ver).String())
		}
	}
	return args
}

func createDefaultTools() map[Name]PkgSrc {
	// go get -u -v xxx
	return map[Name]PkgSrc{
		"gocode":            "github.com/mdempsky/gocode",
		"gopkgs":            "github.com/uudashr/gopkgs/cmd/gopkgs",
		"go-outline":        "github.com/ramya-rao-a/go-outline",
		"go-symbols":        "github.com/acroca/go-symbols",
		"guru":              "golang.org/x/tools/cmd/guru",
		"gorename":          "golang.org/x/tools/cmd/gorename",
		"dlv":               "github.com/go-delve/delve/cmd/dlv",
		"stamblerre_gocode": "github.com/stamblerre/gocode",
		"godef":             "github.com/rogpeppe/godef",
		"goreturns":         "github.com/sqs/goreturns",
		"golint":            "golang.org/x/lint/golint",
		"gotests":           "github.com/cweill/gotests/...",
		"gomodifytags":      "github.com/fatih/gomodifytags",
		"impl":              "github.com/josharian/impl",
		"fillstruct":        "github.com/davidrjenni/reftools/cmd/fillstruct",
		"goplay":            "github.com/haya14busa/goplay/cmd/goplay",
		"godoctor":          " github.com/godoctor/godoctor",
	}
}
