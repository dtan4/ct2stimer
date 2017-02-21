package systemd

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// ServiceData represents data set of systemd Service
type ServiceData struct {
	Name    string
	Command string
	After   string
	User    string
}

// TimerData represents data set of systemd Timer
type TimerData struct {
	Name     string
	Cronspec string
}

// GenerateService generates new systemd Service
func GenerateService(name, command, after, user string) (string, error) {
	body, err := Asset("templates/service.tmpl")
	if err != nil {
		return "", errors.Wrap(err, "failed to load service template")
	}

	tmpl, err := template.New("ct2stimer-service").Parse(string(body))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse service template")
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, &ServiceData{
		Name:    name,
		Command: command,
		After:   after,
		User:    user,
	}); err != nil {
		return "", errors.Wrap(err, "failed to dispatch values in service template")
	}

	return buf.String(), nil
}

// GenerateTimer generates new systemd Timer
func GenerateTimer(name, cronspec string) (string, error) {
	body, err := Asset("templates/timer.tmpl")
	if err != nil {
		return "", errors.Wrap(err, "failed to load timer template")
	}

	tmpl, err := template.New("ct2stimer-timer").Parse(string(body))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse timer template")
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, &TimerData{
		Name:     name,
		Cronspec: cronspec,
	}); err != nil {
		return "", errors.Wrap(err, "failed to dispatch values in timer template")
	}

	return buf.String(), nil
}
