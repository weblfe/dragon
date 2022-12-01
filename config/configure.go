package config

type Configure struct {
	UserScriptsDir []string `json:"scripts_dir" yaml:"scripts_dir"`
	UserTemplatesDir  []string  `json:"templates_dir" yaml:"templates_dir"`
}