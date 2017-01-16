package crontab

import (
	"fmt"
	"strings"
)

// Schedule represents crontab spec and command
type Schedule struct {
	Spec    string
	Command string
}

// Parse parses crontab file and return a list of Schedule
func Parse(crontab string) ([]*Schedule, error) {
	schedules := []*Schedule{}
	lines := strings.Split(crontab, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		ss := strings.SplitN(line, " ", 6)
		if len(ss) < 6 {
			return []*Schedule{}, fmt.Errorf("Invalid format. line: %s", line)
		}

		schedules = append(schedules, &Schedule{
			Spec:    strings.Join(ss[0:5], " "),
			Command: ss[5],
		})
	}

	return schedules, nil
}
