package crontab

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	filename := filepath.Join("..", "testdata", "crontab")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to open testdata. filename: %s, error: %s", filename, err)
	}

	expected := []*Schedule{
		&Schedule{
			Spec:    "0,5,10,15,20,25,30,35,40,45,50,55 * * * *",
			Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task01.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task01 RAILS_ENV=production'",
		},
		&Schedule{
			Spec:    "15 * * * *",
			Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task02.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task02 RAILS_ENV=production'",
		},
		&Schedule{
			Spec:    "30 * * * *",
			Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task04.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task04 RAILS_ENV=production'",
		},
	}

	got, err := Parse(string(body))
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Schedules do not match.\n  expected: %q\n  got:      %q", expected, got)
	}
}

func TestConvertToSystemdCalendar(t *testing.T) {
	testcases := []struct {
		schedule *Schedule
		expected string
	}{
		{
			schedule: &Schedule{
				Spec:    "*/5 * * * *",
				Command: "",
			},
			expected: "*:0,5,10,15,20,25,30,35,40,45,50,55",
		},
		{
			schedule: &Schedule{
				Spec:    "0,5,10,15,20,25,30,35,40,45,50,55 10-12 * * *",
				Command: "",
			},
			expected: "10,11,12:0,5,10,15,20,25,30,35,40,45,50,55", // TODO: 10-12:0,5,...
		},
		{
			schedule: &Schedule{
				Spec:    "0-5 * 1 * *",
				Command: "",
			},
			expected: "*-1 *:0,1,2,3,4,5", // TODO: *:0-5
		},
		{
			schedule: &Schedule{
				Spec:    "23 2,1 * 12 1,6",
				Command: "",
			},
			expected: "Mon,Sat 12-* 1,2:23",
		},
		{
			schedule: &Schedule{
				Spec:    "0,20,40 8-17 * * 1-5",
				Command: "",
			},
			expected: "Mon,Tue,Wed,Thu,Fri 8,9,10,11,12,13,14,15,16,17:0,20,40",
		},
		{
			schedule: &Schedule{
				Spec:    "0 17 * * *",
				Command: "",
			},
			expected: "17:0",
		},
		{
			schedule: &Schedule{
				Spec:    "* * * * 0",
				Command: "",
			},
			expected: "Sun *:*",
		},
		{
			schedule: &Schedule{
				Spec:    "5 * * * *",
				Command: "",
			},
			expected: "*:5",
		},
	}

	for _, tc := range testcases {
		got, err := tc.schedule.ConvertToSystemdCalendar()
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if got != tc.expected {
			t.Errorf("Calendar does not match. expected: %q, actual: %q", tc.expected, got)
		}
	}
}
