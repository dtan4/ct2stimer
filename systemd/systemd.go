package systemd

// go:generate go-bindata -pkg $GOPACKAGE templates/

import (
	"bytes"
	"text/template"
)

// ServiceData represents data set of Systemd Service
type ServiceData struct {
	Name    string
	Command string
}

// GenerateService generates new Systemd Service
func GenerateService(name, command string) (string, error) {
	body, err := Asset("templates/service.tmpl")
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("ct2stimer-service").Parse(string(body))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, &ServiceData{
		Name:    name,
		Command: command,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
