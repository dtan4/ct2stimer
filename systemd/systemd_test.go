package systemd

import (
	"testing"
)

func TestGenerateService(t *testing.T) {
	name := "ct2stimer"
	command := "/bin/bash docker run --rm ubuntu:16.04 echo hello"
	expected := `[Unit]
Dscription=ct2stimer timer service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStart=/bin/bash docker run --rm ubuntu:16.04 echo hello
Type=oneshot
`

	got, err := GenerateService(name, command)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if got != expected {
		t.Errorf("Unit does not match.\n\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}
