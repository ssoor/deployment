package main

import (
	"bytes"
	"text/template"
)

const serviceTmpl = `
FROM {{ .Image }}

ENV {{ .Env }}

WORKDIR {{ .WorkDir }}

COPY .	.

CMD {{ .Command }}
`

func buildDockerfile(image, cmd string, runtime Runtime) (string, error) {
	builder, err := template.New("systemd").Parse(serviceTmpl)
	if nil != err {
		return "", err
	}

	buff := bytes.NewBuffer([]byte{})

	var tmplData = map[string]string{
		"Image":   image,
		"Command": cmd,
		"WorkDir": runtime.WorkDir,
	}

	if err := builder.Execute(buff, tmplData); nil != err {
		return "", err
	}

	return buff.String(), nil
}
