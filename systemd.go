package main

import (
	"bytes"
	"text/template"
)

const serviceTmpl = `
[Unit]
Description={{ .Description }}
Documentation={{ .Documentation }}
Wants=network-online.target
After=network-online.target firewalld.service

[Service]
Type=simple

ExecStart={{ .Command }}
WorkingDirectory={{ .WorkDir }}

# ulimit
LimitCORE=infinity
LimitNPROC=infinity
LimitNOFILE=infinity

# restart the docker process if it exits prematurely
Restart=on-failure
RestartSec=3s
#StartLimitBurst=3
#StartLimitInterval=60s

[Install]
WantedBy=multi-user.target
`

func buildSystemd(name, cmd string, runtime Runtime) (string, error) {
	builder, err := template.New("systemd").Parse(serviceTmpl)
	if nil != err {
		return "", err
	}

	buff := bytes.NewBuffer([]byte{})

	var tmplData = map[string]string{
		"Command":       cmd,
		"WorkDir":       runtime.WorkDir,
		"Description":   "[Deployment] " + name,
		"Documentation": "",
	}

	if err := builder.Execute(buff, tmplData); nil != err {
		return "", err
	}

	return buff.String(), nil
}
