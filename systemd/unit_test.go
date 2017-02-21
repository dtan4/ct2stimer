package systemd

import (
	"testing"
)

func TestGenerateService(t *testing.T) {
	testcases := []struct {
		name     string
		command  string
		after    string
		user     string
		expected string
	}{
		{
			name:    "ct2stimer",
			command: "/bin/bash docker run --rm ubuntu:16.04 echo hello",
			after:   "docker.service",
			user:    "",
			expected: `[Unit]
Description=ct2stimer service unit
After=docker.service

[Service]
TimeoutStartSec=0
ExecStart=/bin/bash docker run --rm ubuntu:16.04 echo hello
Type=oneshot
`,
		},
		{
			name:    "ct2stimer",
			command: "/bin/bash docker run --rm ubuntu:16.04 echo hello",
			after:   "",
			user:    "core",
			expected: `[Unit]
Description=ct2stimer service unit

[Service]
TimeoutStartSec=0
ExecStart=/bin/bash docker run --rm ubuntu:16.04 echo hello
Type=oneshot
User=core
`,
		},
	}

	for _, tc := range testcases {
		got, err := GenerateService(tc.name, tc.command, tc.after, tc.user)
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if got != tc.expected {
			t.Errorf("Service does not match.\n\nexpected:\n%s\n\ngot:\n%s", tc.expected, got)
		}
	}
}

func TestGenerateTimer(t *testing.T) {
	name := "ct2stimer"
	cronspec := "30 * * * *"
	expected := `[Unit]
Description=ct2stimer timer unit

[Timer]
OnCalendar=30 * * * *
`

	got, err := GenerateTimer(name, cronspec)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if got != expected {
		t.Errorf("Timer does not match.\n\nexpected:\n%s\n\ngot:\n%s", expected, got)
	}
}
