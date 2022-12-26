package tools

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"
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
	goProxy := os.Getenv("GOPROXY")
	log.Printf("GOROOT: %s \n", root)
	log.Printf("GOPATH: %s \n", goPath)
	log.Printf("GOPROXY: %s \n", goProxy)
	v, ok := dev.args["GOROOT"]
	if !ok || (root != "" && v == "") {
		dev.args["GOROOT"] = root
	}
	v, ok = dev.args["GOPATH"]
	if !ok || (goPath != "" && v == "") {
		dev.args["GOPATH"] = goPath
	}
}

func (dev *Devops) Exec() (stdOut, stdErr bytes.Buffer, err error) {
	var (
		stdout, stderr bytes.Buffer
		args           = dev.GetArgs()
		name           = dev.GetCommand()
	)
	for _, v := range args {
		stdout.Reset()
		stdErr.Reset()
		if err = dev.GoGet(name, v, &stdout, &stderr); err != nil {
			log.Printf("[ERROR] %v \n", err.Error())
			continue
		}
		install := false
		if stderr.Len() <= 0 {
			install = true
		} else {
			errMsg := stderr.String()
			if strings.Contains(errMsg, "go install") {
				install = true
			} else {
				log.Printf("[ERROR] %v \n", errMsg)
			}
		}
		if install && !strings.Contains(v, "...") {
			stdout.Reset()
			stdErr.Reset()
			err = dev.GoInstall(name, v, &stdout, &stderr)
			errMsg := stderr.String()
			if errMsg != "" {
				log.Printf("[ERROR] %v \n", errMsg)
			}
		}
		if err != nil {
			log.Printf("[ERROR] %v \n", err.Error())
			continue
		}
		if stdout.Len() > 0 {
			log.Printf("%v\n", stdout.String())
		}
	}
	return stdout, stderr, nil
}

func (dev *Devops) GoGet(bin, pkg string, stdout, stderr *bytes.Buffer) error {
	cmd := exec.Command(bin, "get", "-u", "-v", pkg)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	log.Printf(filepath.Base(bin)+" %v %v %v %v \n", "get", "-u", "-v ", pkg)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (dev *Devops) GoInstall(bin, pkg string, stdout, stderr *bytes.Buffer) error {
	cmd := exec.Command(bin, "install", pkg)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	log.Printf(filepath.Base(bin)+" %v %v \n", "install", pkg)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (dev *Devops) GetCommand() string {
	root, ok := dev.args["GOROOT"]
	if !ok || root == "" {
		return "go"
	}
	return dev.getGoBin(root)
}

func (dev *Devops) getGoBin(root string) string {
	var bin = filepath.Join(root, "bin", "go")
	if runtime.GOOS == "windows" {
		bin = filepath.Join(root, "bin", "go.exe")
	}
	return bin
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
		sort.Strings(args)
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
	sort.Strings(args)
	return args
}

func (dev *Devops) Lists() {
	w := new(tabwriter.Writer)
	var (
		max   int
		names []string
	)
	for v, s := range dev.tools {
		name := v.String()
		names = append(names, name)
		if size := len(name) + len(s.String()); max < size {
			max = size
		}
	}
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, max+8, 8, '\t', 0)
	_, _ = w.Write([]byte("name\tpkg_src\t\n"))
	_, _ = w.Write([]byte(strings.Repeat("-", max+8) + "\n"))
	sort.Strings(names)
	for _, n := range names {
		_, _ = w.Write([]byte(fmt.Sprintf("%s\t%s\t\n", n, dev.tools[Name(n)])))
	}

	_ = w.Flush()
}

func createDefaultTools() map[Name]PkgSrc {
	// go get -u -v xxx
	return map[Name]PkgSrc{
		"gocode":            "github.com/mdempsky/gocode",
		"gopkgs":            "github.com/tpng/gopkgs",
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
		"godoctor":          "github.com/godoctor/godoctor",
	}
}
