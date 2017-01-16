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
