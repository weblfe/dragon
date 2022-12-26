package render

import (
	"errors"
	"github.com/shiyanhui/hero"
)

const (
	extension = ".tgo"
)

type Render interface {
	Generate(arg *GenerateArg) error
}

type GenerateArg struct {
	Source     string
	Dest       string
	PkgName    string
	Extensions []string
}

type renderImpl struct {
}

func (r *renderImpl) Generate(arg *GenerateArg) error {
	arg.Extensions = append(arg.Extensions, extension)
	if !hero.CheckExtension(arg.Source, arg.Extensions) {
		return errors.New(`extension miss`)
	}
	if arg.PkgName == "" {
		arg.PkgName = "template"
	}
	hero.Generate(arg.Source, arg.Dest, arg.PkgName, arg.Extensions)
	return nil
}

func NewRender() Render {
	return &renderImpl{}
}
