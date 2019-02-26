package main

type DeployFile struct {
	Targets map[string]FileTarget `json:"deploys" yaml:"deploys"`
}

type FileTarget struct {
	Target  `json:",inline" yaml:",inline"`
	Imports []string `json:"imports,omitempty" yaml:"imports"`
}

type Target struct {
	Command string    `json:"cmd,omitempty" yaml:"cmd"`
	Env     []Env     `json:"env,omitempty" yaml:"env"`
	Props   []Port    `json:"ports,omitempty" yaml:"ports"`
	Assets  []Assets  `json:"assets,omitempty" yaml:"assets"`
	Runtime []Runtime `json:"runtime,omitempty" yaml:"runtime"`
}

type Env struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type Assets struct {
	Name   string `json:"name" yaml:"name"`
	Target string `json:"target" yaml:"target"`
	Source string `json:"source" yaml:"source"`
}

type Port struct {
	Name          string `json:"name" yaml:"name"`
	Protocol      string `json:"protocol" yaml:"protocol"`
	ContainerPort string `json:"containerPort" yaml:"containerPort"`
}

type Runtime struct {
	Type    string `json:"type" yaml:"type"`
	SSH     string `json:"ssh,omitempty" yaml:"ssh"`
	Image   string `json:"image,omitempty" yaml:"image"`
	WorkDir string `json:"workDir" yaml:"workDir"`
}
